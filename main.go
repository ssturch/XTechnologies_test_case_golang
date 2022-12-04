package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Главная функция отвечающая за получение, переработку и отправку данных
func mainprocess(pgdb *sql.DB, i int, j int) (int, int) {
	var rawResCBR ResultXML
	var rawResBTC map[string]interface{}
	var resCBR map[string]map[string]interface{}
	var resBTC map[string]interface{}

	err = nil
	rawResCBR, err = getXMLfromCBR()
	if err != nil {
		errorLogger(err)
		return i, j
	}
	rawResBTC, err = getJSONfromKUCOIN()
	if err != nil {
		errorLogger(err)
		return i, j
	}
	resCBR, err = remakeXMLfromCBR(rawResCBR)
	if err != nil {
		errorLogger(err)
		return i, j
	}
	resBTC, err = remakeJSONfromKUCOIN(rawResBTC)
	if err != nil {
		errorLogger(err)
		return i, j
	}
	RUBBTCValue := calcRUBBTC(resCBR["USD"]["value"].(float64), resBTC["value"].(float64))
	if comparseDataBy_btcusd_btcrub(pgdb, resBTC["value"].(float64), RUBBTCValue) == true {
		_, err = pgdb.Exec(addDataTo_btcusd_btcrub(), resBTC["time"], resBTC["value"], RUBBTCValue)
		if err != nil {
			errorLogger(err)
			return i, j
		}
	}
	mapValToBTC := calcAllValutesBTC(resCBR, resBTC["value"].(float64))
	if comparseDataBy_cbr_btc(pgdb, mapValToBTC, i) == true {
		var tempTime int
		for el, val := range resCBR {
			tempTime = resBTC["time"].(int)
			_, err = pgdb.Exec(addDataTo_cbr_btc(), i, tempTime, el, val["name"], mapValToBTC[el])
			if err != nil {
				errorLogger(err)
				return i, j
			}
		}
		_, err = pgdb.Exec(addDataTo_cbr_btc(), i, tempTime, "RUB", "Российский рубль", RUBBTCValue)
		i++
	}
	if comparseDataBy_cbr_rub(pgdb, resCBR, j) == true {
		var tempTime int64
		for el, val := range resCBR {
			tempTime = val["time"].(int64)
			_, err = pgdb.Exec(addDataTo_cbr_rub(), j, tempTime, el, val["name"], val["value"])
			if err != nil {
				errorLogger(err)
				return i, j
			}
		}
		j++

	}
	APIResult = make(map[string]interface{})
	APIResult["LastValue_btcusd"], err = getLastValue_btcusd(pgdb)
	APIResult["History_btcusd"], err = getHistory_btcusd(pgdb)
	APIResult["LastValues_cbr_rub"], err = getLastValue_cbr_rub(pgdb)
	APIResult["History_cbr_rub"], err = getHistory_cbr_rub(pgdb)
	APIResult["LastValues_btc_curr"], err = getLastValue_btc_curr(pgdb)
	APIResult["History_btc_curr"], err = getHistory_btc_curr(pgdb)
	return i, j
}

// Горутина основного процесса (вызывает mainprocess по счетчику)
func worker(pgdb *sql.DB) {
	i := getStartIdBy_cbr_btc(pgdb)
	j := getStartIdBy_cbr_rub(pgdb)
	heartbeatFirst := time.After(100 * time.Millisecond)
	heartbeat := time.Tick(10 * time.Second)
	for {
		select {
		case <-heartbeatFirst:
			i, j = mainprocess(pgdb, i, j)
		case <-heartbeat:
			i, j = mainprocess(pgdb, i, j)
		}
	}
}

var APIResult map[string]interface{}
var err error

func main() {
	var pgdb *sql.DB
	pgdb, err = pgdbconnect()
	if err != nil {
		errorLogger(err)
		os.Exit(0)
	}
	tempQueryVar := createTables()
	for i := range tempQueryVar {
		_, err = pgdb.Exec(tempQueryVar[i])
		if err != nil {
			errorLogger(err)
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
				pgdb.Close()
				exitChan <- 0
			case os.Interrupt:
				fmt.Println("Catch: SIGNAL INTERRUPT | Server stopped | DB Closed")
				pgdb.Close()
				exitChan <- 0
			case syscall.SIGTERM:
				fmt.Println("Catch: SIGNAL TERMINATE | Server stopped | DB Closed")
				pgdb.Close()
				exitChan <- 0
			case syscall.SIGKILL:
				fmt.Println("Catch: SIGNAL KILL | Server stopped | DB Closed")
				pgdb.Close()
				exitChan <- 0
			}
		}
	}()
	// Отслеживание входящих запросов по REST API
	go server()
	// Горутина основного процесса (вызывает mainprocess по счетчику)
	go worker(pgdb)
	exitCode := <-exitChan
	os.Exit(exitCode)
}
