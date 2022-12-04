package main

import (
	"net/http"
	"os"
)

const (
	//Host = "localhost"
	Port = "8080"
)

// Отслеживание endpoint'ов
func server() {

	http.HandleFunc("/api/btcusdt", btcusdtHandler)
	http.HandleFunc("/api/currencies", currenciesRUBHandler)
	http.HandleFunc("/api/latest", currenciesBTCHandler)
	http.HandleFunc("/api/latest/", currenciesBTCbyCHARHandler)
	realhost, _ := os.Hostname()
	err = http.ListenAndServe(realhost+":"+Port, nil)
	if err != nil {
		errorLogger(err)
	}
}
