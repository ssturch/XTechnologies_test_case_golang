package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Создание табилц необходимых для работы программы
func createTables() []string {
	query1 := "CREATE TABLE IF NOT EXISTS btcusd_btcrub (" +
		"id SERIAL PRIMARY KEY," +
		"time_UTC bigint NOT NULL UNIQUE," +
		"btcusd_value double precision NOT NULL," +
		"btcrub_value double precision NOT NULL);"
	query2 := "CREATE TABLE IF NOT EXISTS cbr_btc (" +
		"id SERIAL PRIMARY KEY," +
		"count integer NOT NULL," +
		"time_UTC bigint NOT NULL," +
		"char_code text NOT NULL," +
		"valute_name text NOT NULL," +
		"value_btc double precision NOT NULL);"
	query3 := "CREATE TABLE IF NOT EXISTS cbr_rub (" +
		"id SERIAL PRIMARY KEY," +
		"count integer NOT NULL," +
		"time_UTC bigint NOT NULL," +
		"char_code text NOT NULL," +
		"valute_name text NOT NULL," +
		"value_rub double precision NOT NULL);"

	queryArr := []string{query1, query2, query3}
	return queryArr
}

// Шаблон запроса INSERT
func addDataTo_btcusd_btcrub() string {
	query := "INSERT INTO btcusd_btcrub(time_UTC, btcusd_value, btcrub_value)" +
		"VALUES ($1, $2, $3);"
	return query
}

// Шаблон запроса INSERT
func addDataTo_cbr_btc() string {
	query := "INSERT INTO cbr_btc(count, time_UTC, char_code, valute_name, value_btc)" +
		"VALUES ($1, $2, $3, $4, $5);"
	return query
}

// Шаблон запроса INSERT
func addDataTo_cbr_rub() string {
	query := "INSERT INTO cbr_rub(count, time_UTC, char_code, valute_name, value_rub)" +
		"VALUES ($1, $2, $3, $4, $5);"
	return query
}

// Сравнение старых данных из базы с новыми из запроса GET к серверам (возвращает true если данные отличаются)
func comparseDataBy_btcusd_btcrub(pgdb *sql.DB, resBTC float64, RUBBTCValue float64) bool {
	var temp_resBTC float64
	var temp_RUBBTCValue float64
	query := "SELECT btcusd_value, btcrub_value FROM btcusd_btcrub ORDER BY time_utc DESC LIMIT 1;"
	resQuery := pgdb.QueryRow(query)
	err = resQuery.Scan(&temp_resBTC, &temp_RUBBTCValue)
	if err != nil {
		return true
	}
	if resBTC != temp_resBTC {
		return true
	}
	if RUBBTCValue != temp_RUBBTCValue {
		return true
	}
	return false
}

// Сравнение старых данных из базы с новыми из запроса GET к серверам (возвращает true если данные отличаются)
func comparseDataBy_cbr_btc(pgdb *sql.DB, mapValToBTC map[string]float64, i int) bool {
	var resQuery *sql.Rows
	if i == 1 {
		return true
	}
	query := "SELECT char_code, value_btc FROM cbr_btc WHERE count = $1;"
	resQuery, err = pgdb.Query(query, i-1)
	if err != nil {
		return true
	}
	for resQuery.Next() {
		var temp_char_code string
		var temp_value_btc float64
		err = resQuery.Scan(&temp_char_code, &temp_value_btc)
		if mapValToBTC[temp_char_code] != temp_value_btc {
			resQuery.Close()
			return true
		}
	}
	return false
}

// Сравнение старых данных из базы с новыми из запроса GET к серверам (возвращает true если данные отличаются)
func comparseDataBy_cbr_rub(pgdb *sql.DB, mapValRub map[string]map[string]interface{}, j int) bool {
	var resQuery *sql.Rows
	if j == 1 {
		return true
	}
	query := "SELECT char_code, value_rub FROM cbr_rub WHERE count = $1;"
	resQuery, err = pgdb.Query(query, j-1)
	if err != nil {
		fmt.Println(err)
		return true
	}
	for resQuery.Next() {
		var temp_char_code string
		var temp_value_rub float64
		err = resQuery.Scan(&temp_char_code, &temp_value_rub)
		if mapValRub[temp_char_code]["value"] != temp_value_rub {
			fmt.Println(temp_char_code)
			fmt.Println(mapValRub[temp_char_code])
			fmt.Println(temp_value_rub)
			resQuery.Close()
			return true
		}
	}
	return false
}

// Получение стартового ID, в случае если сервис был перезапущен, или используется старая БД, возвращает INT
func getStartIdBy_cbr_btc(pgdb *sql.DB) int {
	var startDigit int
	query := "SELECT count FROM cbr_btc ORDER BY id DESC LIMIT 1;"
	resQuery := pgdb.QueryRow(query)
	err = resQuery.Scan(&startDigit)
	if err != nil {
		return 1
	}
	return startDigit + 1
}

// Получение стартового ID, в случае если сервис был перезапущен, или используется старая БД, возвращает INT
func getStartIdBy_cbr_rub(pgdb *sql.DB) int {
	var startDigit int
	query := "SELECT count FROM cbr_rub ORDER BY id DESC LIMIT 1;"
	resQuery := pgdb.QueryRow(query)
	err = resQuery.Scan(&startDigit)
	if err != nil {
		return 1
	}
	return startDigit + 1
}

// Получение последней строчки из БД запросом SELECT, возвращает JSON
func getLastValue_btcusd(pgdb *sql.DB) (float64, error) {
	var result float64
	query := "SELECT value_btc FROM cbr_btc WHERE char_code = 'USD' ORDER BY count DESC LIMIT 1;"
	resQuery := pgdb.QueryRow(query)
	err = resQuery.Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Получение всех данных из БД запросом SELECT, возвращает JSON
func getHistory_btcusd(pgdb *sql.DB) ([]byte, error) {
	type JSONResult struct {
		Total   int
		History []map[string]interface{}
	}
	var resQuery *sql.Rows
	var tempSlice []map[string]interface{}
	var result []byte
	query := "SELECT time_utc, value_btc FROM cbr_btc WHERE char_code = 'USD';"
	resQuery, err = pgdb.Query(query)
	if err != nil {
		return nil, err
	}
	i := 0
	for resQuery.Next() {
		var tempTime int
		var tempValue float64
		err = resQuery.Scan(&tempTime, &tempValue)
		tempMap := make(map[string]interface{}, 2)
		tempMap["timestamp"] = tempTime
		tempMap["value"] = tempValue
		tempSlice = append(tempSlice, tempMap)
		i++
	}
	resQuery.Close()
	tempResult := JSONResult{Total: i, History: tempSlice}

	result, err = json.Marshal(tempResult)
	if err != nil {
		return nil, err
	}
	return result, err
}

// Получение последней строчки из БД запросом SELECT, возвращает JSON
func getLastValue_cbr_rub(pgdb *sql.DB) ([]byte, error) {
	type JSONResult struct {
		Values map[string]float64
	}
	var resQuery *sql.Rows
	var tempTime int
	var result []byte

	query := "SELECT time_utc, char_code, value_rub FROM cbr_rub WHERE count = (SELECT count FROM cbr_rub ORDER BY count DESC LIMIT 1)"
	resQuery, err = pgdb.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tempMap := make(map[string]float64)
	for resQuery.Next() {
		var tempChar string
		var tempValue float64
		err = resQuery.Scan(&tempTime, &tempChar, &tempValue)
		tempMap[tempChar] = tempValue
	}
	resQuery.Close()
	tempResult := JSONResult{Values: tempMap}
	result, err = json.Marshal(tempResult)
	if err != nil {
		return nil, err
	}
	return result, err
}

// Получение всех данных из БД запросом SELECT, возвращает JSON
func getHistory_cbr_rub(pgdb *sql.DB) ([]byte, error) {
	type JSONResult struct {
		Total   int
		History []map[string]interface{}
	}
	var resQuery *sql.Rows
	var ok bool
	var result []byte

	query := "SELECT count, time_utc, char_code, value_rub FROM cbr_rub ORDER BY count;"
	resQuery, err = pgdb.Query(query)
	if err != nil {
		return nil, err
	}

	i := 0
	oldTempCount := 0
	var tempSlice []map[string]interface{}
	tempMap := make(map[string]interface{})

	for ok = resQuery.Next(); ok; {
		var tempCount int
		var tempChar string
		var tempValue float64
		var tempTime int
		err = resQuery.Scan(&tempCount, &tempTime, &tempChar, &tempValue)
		if resQuery.Next() == false {
			tempMap[tempChar] = tempValue
			tempTimeStr := time.Unix(int64(tempTime), 0).Format("2006-01-02")
			tempMap["date"] = tempTimeStr
			tempSlice = append(tempSlice, tempMap)
			i++
			break
		}
		if tempCount != oldTempCount && oldTempCount != 0 {
			oldTempCount = tempCount
			tempMap[tempChar] = tempValue
			tempTimeStr := time.Unix(int64(tempTime), 0).Format("2006-01-02")
			tempMap["date"] = tempTimeStr
			tempSlice = append(tempSlice, tempMap)
			tempMap = make(map[string]interface{})
			i++
		} else {
			oldTempCount = tempCount
			tempMap[tempChar] = tempValue
		}
	}
	resQuery.Close()
	tempResult := JSONResult{Total: i, History: tempSlice}

	result, err = json.Marshal(tempResult)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}

// Получение последней строчки из БД запросом SELECT, возвращает JSON
func getLastValue_btc_curr(pgdb *sql.DB) ([]byte, error) {
	type JSONResult struct {
		//Date   string
		Values map[string]float64
	}
	var resQuery *sql.Rows
	var tempTime int
	var result []byte

	query := "SELECT time_utc, char_code, value_btc FROM cbr_btc WHERE count = (SELECT count FROM cbr_btc ORDER BY count DESC LIMIT 1)"
	resQuery, err = pgdb.Query(query)
	if err != nil {
		return nil, err
	}

	tempMap := make(map[string]float64)
	for resQuery.Next() {
		var tempChar string
		var tempValue float64
		err = resQuery.Scan(&tempTime, &tempChar, &tempValue)
		tempMap[tempChar] = tempValue
	}
	resQuery.Close()
	//intToUNIXTime := time.Unix(int64(tempTime), 0)
	//convertedUNIXTime := intToUNIXTime.Format("2006-01-02")
	//tempResult := JSONResult{Date: convertedUNIXTime, Values: tempMap}
	tempResult := JSONResult{Values: tempMap}
	result, err = json.Marshal(tempResult)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}

// Получение всех данных из БД запросом SELECT, возвращает JSON
func getHistory_btc_curr(pgdb *sql.DB) ([]byte, error) {
	type JSONResult struct {
		Total   int
		History []map[string]interface{}
	}
	var resQuery *sql.Rows
	var ok bool
	var result []byte

	query := "SELECT count, time_utc, char_code, value_btc FROM cbr_btc ORDER BY count;"
	resQuery, err = pgdb.Query(query)
	if err != nil {
		return nil, err
	}

	i := 0
	oldTempCount := 0
	var tempSlice []map[string]interface{}
	tempMap := make(map[string]interface{})

	for ok = resQuery.Next(); ok; {
		var tempCount int
		var tempChar string
		var tempValue float64
		var tempTime int
		err = resQuery.Scan(&tempCount, &tempTime, &tempChar, &tempValue)
		if resQuery.Next() == false {
			tempMap[tempChar] = tempValue
			tempMap["timestamp"] = tempTime
			tempSlice = append(tempSlice, tempMap)
			i++
			break
		}
		if tempCount != oldTempCount && oldTempCount != 0 {
			oldTempCount = tempCount
			tempMap[tempChar] = tempValue
			tempMap["timestamp"] = tempTime
			tempSlice = append(tempSlice, tempMap)
			tempMap = make(map[string]interface{})
			i++
		} else {
			oldTempCount = tempCount
			tempMap[tempChar] = tempValue
		}
	}
	resQuery.Close()
	tempResult := JSONResult{Total: i, History: tempSlice}

	result, err = json.Marshal(tempResult)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, err
}
