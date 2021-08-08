package main

import (
	"flag"
	"fmt"
	"github.com/comhttp/explorer/exp"
	"net/http"
	"time"
)

func main() {
	// Get cmd line parameters
	coin := flag.String("coin", "coin", "Coin")
	bind := flag.String("bind", "localhost:15500", "HTTP server bind in format addr:port")
	loglevel := flag.String("loglevel", "info", "Logging level (debug, info, warn, error)")
	flag.Parse()

	e := exp.NewJORMexplorer(*loglevel, *coin)

	ticker := time.NewTicker(12 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				exp.ExploreCoin(e, *coin)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	fmt.Println("JORM explorer is listening: ", *bind)
	// Start HTTP server
	http.ListenAndServe(*bind, nil)
}
