package web

import (
	oth "XTapi/internal"
	"context"
	"net/http"
)

const (
	Host = "localhost"
	Port = "8080"
)

// Отслеживание endpoint'ов
func Server(ctx context.Context, apires *map[string]interface{}) {

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				var err error
				http.HandleFunc("/api/btcusdt", btcusdtHandler(apires))
				http.HandleFunc("/api/currencies", currenciesRUBHandler(apires))
				http.HandleFunc("/api/latest", currenciesBTCHandler(apires))
				http.HandleFunc("/api/latest/", currenciesBTCbyCHARHandler(apires))
				//realhost, _ := os.Hostname()
				realhost := Host // для отладки
				err = http.ListenAndServe(realhost+":"+Port, nil)
				//err = http.ListenAndServe(realhost, nil)
				if err != nil {
					oth.ErrorLogger(&oth.CustomErr{
						Tp:    "Internal",
						Cause: "http.ListenAndServe",
						Text:  "",
						Err:   err,
					})
				}
			}
		}
	}()
}
