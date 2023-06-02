package web

import (
	"context"
	"fmt"
	"net/http"
)

const (
	Host = "localhost"
	Port = "8080"
)

// Отслеживание endpoint'ов
func Server(ctx context.Context, apires *map[string]interface{}) {

	var err error
	srv := &http.Server{Addr: ":" + Port}
	http.HandleFunc("/api/btcusdt", btcusdtHandler(apires))
	http.HandleFunc("/api/currencies", currenciesRUBHandler(apires))
	http.HandleFunc("/api/latest", currenciesBTCHandler(apires))
	http.HandleFunc("/api/latest/", currenciesBTCbyCHARHandler(apires))

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			fmt.Println("Listen and Serve process stopped!")
			return
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)
	fmt.Println("Server stopped!")

}
