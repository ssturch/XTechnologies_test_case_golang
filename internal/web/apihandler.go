package web

import (
	oth "XTapi/internal"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Хэндлер для вывода результатов запроса на /api/btcusdt
func btcusdtHandler(apires *map[string]interface{}) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		var tempValue float64
		var tempResultJSON []byte
		var err error
		if r.Method == "GET" {
			if *apires == nil {
				oth.ErrorLogger(errors.New("Data cache is nil, please try again later"))
				return
			}
			tempValue = (*apires)["LastValue_btcusd"].(float64)

			if err != nil {
				oth.ErrorLogger(err)
				fmt.Fprintf(w, err.Error())
			} else {

				fmt.Fprintf(w, strconv.FormatFloat(tempValue, 'g', -1, 64))
				//btcusdtView(r.Method)
			}
		}
		if r.Method == "POST" {
			tempResultJSON = (*apires)["History_btcusd"].([]byte)
			if err != nil {
				oth.ErrorLogger(err)
				fmt.Fprintf(w, err.Error())
			} else {
				res := string(RemakeJSONtoSample(tempResultJSON))
				fmt.Fprintf(w, res)
			}

		}
	}
}

// Хэндлер для вывода результатов запроса на /api/currencies
func currenciesRUBHandler(apires *map[string]interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tempResultJSON []byte

		if r.Method == "GET" {
			if *apires == nil {
				fmt.Fprintf(w, "Data Cache is nil, please try again later")
				return
			}
			tempResultJSON = (*apires)["LastValues_cbr_rub"].([]byte)
			res := string(RemakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
			return
		}

		if r.Method == "POST" {
			tempResultJSON = (*apires)["History_cbr_rub"].([]byte)
			res := string(RemakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)
			return
		}
	}
}

// Хэндлер для вывода результатов запроса на /api/latest
func currenciesBTCHandler(apires *map[string]interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var tempResultJSON []byte
		if *apires == nil {
			fmt.Fprintf(w, "Data Cache is nil, please try again later")
			return
		}
		if r.Method == "GET" {
			tempResultJSON = (*apires)["LastValues_btc_curr"].([]byte)
			res := string(RemakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)

		}
		if r.Method == "POST" {
			tempResultJSON = (*apires)["History_btc_curr"].([]byte)
			res := string(RemakeJSONtoSample(tempResultJSON))
			fmt.Fprintf(w, res)

		}
	}
}

// Хэндлер для вывода результатов запроса на /api/latest/{char}
func currenciesBTCbyCHARHandler(apires *map[string]interface{}) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		charcode := strings.Replace(r.URL.RequestURI(), "/api/latest/", "", -1)
		var tempResultJSON []byte
		if r.Method == "GET" {
			type JSONResult struct {
				Values map[string]float64
			}
			if *apires == nil {
				fmt.Fprintf(w, "Data Cache is nil, please try again later")
				return
			}
			var res JSONResult
			tempResultJSON = (*apires)["LastValues_btc_curr"].([]byte)
			err := json.Unmarshal(tempResultJSON, &res)
			if err != nil {
				oth.ErrorLogger(&oth.DecodeUnmarshallError{
					Do:    "Unmarshalling",
					Cause: "JSON from map",
					Err:   err,
				})
				fmt.Fprintf(w, err.Error())
				return
			} else {
				result := res.Values[charcode]
				fmt.Fprintf(w, strconv.FormatFloat(result, 'f', -1, 64))
				return
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
			tempResultJSON = (*apires)["History_btc_curr"].([]byte)
			err := json.Unmarshal(tempResultJSON, &oldvalues)
			if err != nil {
				oth.ErrorLogger(&oth.DecodeUnmarshallError{
					Do:    "Unmarshalling",
					Cause: "JSON from map",
					Err:   err,
				})
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
				oth.ErrorLogger(&oth.DecodeUnmarshallError{
					Do:    "Unmarshalling",
					Cause: "JSON from map",
					Err:   err,
				})
				fmt.Fprintf(w, err.Error())
				return
			} else {
				res := string(RemakeJSONtoSample(resultJSON))
				fmt.Fprintf(w, res)
				return
			}
		}
	}
}
