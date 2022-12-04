package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Хэндлер для вывода результатов запроса на /api/btcusdt
func btcusdtHandler(w http.ResponseWriter, r *http.Request) {
	var tempValue float64
	var tempResultJSON []byte
	if r.Method == "GET" {
		tempValue = APIResult["LastValue_btcusd"].(float64)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {

			fmt.Fprintf(w, strconv.FormatFloat(tempValue, 'g', -1, 64))
			//btcusdtView(r.Method)
		}
	}
	if r.Method == "POST" {
		tempResultJSON = APIResult["History_btcusd"].([]byte)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			res := string(remakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
			//btcusdtView(r.Method)
		}

	}
}

// Хэндлер для вывода результатов запроса на /api/currencies
func currenciesRUBHandler(w http.ResponseWriter, r *http.Request) {
	var tempResultJSON []byte
	if r.Method == "GET" {
		tempResultJSON = APIResult["LastValues_cbr_rub"].([]byte)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			res := string(remakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
		}
	}
	if r.Method == "POST" {
		tempResultJSON = APIResult["History_cbr_rub"].([]byte)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			res := string(remakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
		}
	}
}

// Хэндлер для вывода результатов запроса на /api/latest
func currenciesBTCHandler(w http.ResponseWriter, r *http.Request) {
	var tempResultJSON []byte
	if r.Method == "GET" {
		tempResultJSON = APIResult["LastValues_btc_curr"].([]byte)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			res := string(remakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
		}
	}
	if r.Method == "POST" {
		tempResultJSON = APIResult["History_btc_curr"].([]byte)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			res := string(remakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
		}
	}
}

// Хэндлер для вывода результатов запроса на /api/latest/{char}
func currenciesBTCbyCHARHandler(w http.ResponseWriter, r *http.Request) {
	charcode := strings.Replace(r.URL.RequestURI(), "/api/latest/", "", -1)
	var tempResultJSON []byte
	if r.Method == "GET" {
		type JSONResult struct {
			Values map[string]float64
		}
		var res JSONResult
		tempResultJSON = APIResult["LastValues_btc_curr"].([]byte)
		err = json.Unmarshal(tempResultJSON, &res)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			result := res.Values[charcode]
			fmt.Fprintf(w, strconv.FormatFloat(result, 'f', -1, 64))
		}
	}
	if r.Method == "POST" {
		type JSONResult struct {
			Total   int
			History []map[string]interface{}
		}
		var oldvalues JSONResult
		var newvalues JSONResult
		var newhistory []map[string]interface{}
		var resultJSON []byte
		tempResultJSON = APIResult["History_btc_curr"].([]byte)
		err = json.Unmarshal(tempResultJSON, &oldvalues)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
			return
		}
		newvalues.Total = oldvalues.Total
		for _, value := range oldvalues.History {
			if value[charcode] == nil {
				continue
			} else {
				newvalue := make(map[string]interface{}, 2)
				newvalue["timestamp"] = value["timestamp"].(float64)
				newvalue[charcode] = value[charcode].(float64)
				newhistory = append(newhistory, newvalue)
			}
		}
		newvalues.History = newhistory
		resultJSON, err = json.Marshal(newvalues)
		if err != nil {
			errorLogger(err)
			fmt.Fprintf(w, err.Error())
		} else {
			res := string(remakeJSONtoSample(resultJSON))
			fmt.Fprintf(w, res)
		}
	}
}
