package main

import (
	oth "XTapi/internal"
	clc "XTapi/internal/calc"
	dbp "XTapi/internal/db"
	wb "XTapi/internal/web"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	cbrURL    = "http://www.cbr.ru/scripts/XML_daily.asp"
	kukoinURL = "https://api.kucoin.com/api/v1/market/stats?symbol=BTC-USDT"
)

// Главная функция отвечающая за получение, переработку и отправку данных
func mainprocess(pgdb *sql.DB, apires *map[string]interface{}, mut *sync.Mutex, i int, j int) (int, int) {
	var rawResCBR wb.ResultXML
	var rawResBTC map[string]interface{}
	var resCBR map[string]map[string]interface{}
	var resBTC map[string]interface{}
	var err error

	rawResCBR, err = wb.GetXMLfromCBR(cbrURL)
	if err != nil {
		oth.ErrorLogger(err)
		return i, j
	}
	rawResBTC, err = wb.GetJSONfromKUCOIN(kukoinURL)
	if err != nil {
		oth.ErrorLogger(err)
		return i, j
	}
	resCBR, err = wb.RemakeXMLfromCBR(rawResCBR)
	if err != nil {
		oth.ErrorLogger(err)
		return i, j
	}
	resBTC, err = wb.RemakeJSONfromKUCOIN(rawResBTC)
	if err != nil {
		oth.ErrorLogger(err)
		return i, j
	}
	RUBBTCValue := clc.ClcRUBBTC(resCBR["USD"]["value"].(float64), resBTC["value"].(float64))
	if dbp.ComparseDataBy_btcusd_btcrub(pgdb, resBTC["value"].(float64), RUBBTCValue) == true {
		_, err = pgdb.Exec(dbp.AddDataTo_btcusd_btcrub(), resBTC["time"], resBTC["value"], RUBBTCValue)
		if err != nil {
			oth.ErrorLogger(&oth.DBError{
				QueryFunc: "AddDataTo_btcusd_btcrub()",
				Err:       err,
			})
			return i, j
		}
	}
	mapValToBTC := clc.ClcAllValutesBTC(resCBR, resBTC["value"].(float64))
	if dbp.ComparseDataBy_cbr_btc(pgdb, mapValToBTC, i) == true {
		var tempTime int
		for el, val := range resCBR {
			tempTime = resBTC["time"].(int)
			_, err = pgdb.Exec(dbp.AddDataTo_cbr_btc(), i, tempTime, el, val["name"], mapValToBTC[el])
			if err != nil {
				oth.ErrorLogger(&oth.DBError{
					QueryFunc: "AddDataTo_cbr_btc()",
					Err:       err,
				})
				return i, j
			}
		}
		pgdb.Exec(dbp.AddDataTo_cbr_btc(), i, tempTime, "RUB", "Российский рубль", RUBBTCValue)
		i++
	}
	if dbp.ComparseDataBy_cbr_rub(pgdb, resCBR, j) == true {
		var tempTime int64
		for el, val := range resCBR {
			tempTime = val["time"].(int64)
			_, err = pgdb.Exec(dbp.AddDataTo_cbr_rub(), j, tempTime, el, val["name"], val["value"])
			if err != nil {
				oth.ErrorLogger(&oth.DBError{
					QueryFunc: "AddDataTo_cbr_btc()",
					Err:       err,
				})
				return i, j
			}
		}
		j++
	}
	mut.Lock()
	err = nil
	(*apires)["LastValue_btcusd"], err = dbp.GetLastValue_btcusd(pgdb)
	if err != nil {
		oth.ErrorLogger(err)
	}
	(*apires)["History_btcusd"], err = dbp.GetHistory_btcusd(pgdb)
	if err != nil {
		oth.ErrorLogger(err)
	}
	(*apires)["LastValues_cbr_rub"], err = dbp.GetLastValue_cbr_rub(pgdb)
	if err != nil {
		oth.ErrorLogger(err)
	}
	(*apires)["History_cbr_rub"], err = dbp.GetHistory_cbr_rub(pgdb)
	if err != nil {
		oth.ErrorLogger(err)
	}
	(*apires)["LastValues_btc_curr"], err = dbp.GetLastValue_btc_curr(pgdb)
	if err != nil {
		oth.ErrorLogger(err)
	}
	(*apires)["History_btc_curr"], err = dbp.GetHistory_btc_curr(pgdb)
	if err != nil {
		oth.ErrorLogger(err)
	}
	mut.Unlock()
	return i, j
}

// Горутина основного процесса (вызывает mainprocess по счетчику)
func Worker(ctx context.Context, pgdb *sql.DB, apires *map[string]interface{}) {
	var mutex sync.Mutex
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				i := dbp.GetStartIdBy_cbr_btc(pgdb)
				j := dbp.GetStartIdBy_cbr_rub(pgdb)
				heartbeatFirst := time.After(100 * time.Millisecond)
				heartbeat := time.Tick(10 * time.Second)
				for {
					select {
					case <-heartbeatFirst:
						i, j = mainprocess(pgdb, apires, &mutex, i, j)
					case <-heartbeat:
						i, j = mainprocess(pgdb, apires, &mutex, i, j)
					}
				}
			}
		}
	}()
}

func main() {

	var pgdb *sql.DB
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	APIResult := make(map[string]interface{})
	pgdb, err = dbp.Pgdbconnect()
	if err != nil {
		oth.ErrorLogger(&oth.DBError{
			QueryFunc: "Connect to DB",
			Err:       err,
		})
		os.Exit(0)
	}
	tempQueryVar := dbp.CreateTables()
	for i := range tempQueryVar {
		_, err = pgdb.Exec(tempQueryVar[i])
		if err != nil {
			oth.ErrorLogger(&oth.DBError{
				QueryFunc: "CreateTables()",
				Err:       err,
			})
			break
			os.Exit(0)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	exitChan := make(chan int)

	// Отслеживание системных сигналов (частично совместим с Windows)
	go func() {
		for {
			s := <-sigChan
			switch s {
			case syscall.SIGINT:
				fmt.Println("Catch: SIGNAL INTERRUPT | Server stopped | DB Closed")
				cancel()
				pgdb.Close()
				exitChan <- 0
			case os.Interrupt:
				fmt.Println("Catch: SIGNAL INTERRUPT | Server stopped | DB Closed")
				cancel()
				pgdb.Close()
				exitChan <- 0
			case syscall.SIGTERM:
				fmt.Println("Catch: SIGNAL TERMINATE | Server stopped | DB Closed")
				cancel()
				pgdb.Close()
				exitChan <- 0
			case syscall.SIGKILL:
				fmt.Println("Catch: SIGNAL KILL | Server stopped | DB Closed")
				cancel()
				pgdb.Close()
				exitChan <- 0
			}
		}
	}()

	// Отслеживание входящих запросов по REST API
	wb.Server(ctx, &APIResult)
	// Горутина основного процесса (вызывает mainprocess по счетчику)
	Worker(ctx, pgdb, &APIResult)
	exitCode := <-exitChan
	os.Exit(exitCode)
}
