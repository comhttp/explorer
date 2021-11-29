package app

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

//type Explorer struct {
//	Status map[string]*BlockchainStatus `json:"status"`
//	Status map[string]*BlockchainStatus `json:"status"`
//}

func MainLoop(path, command, coin string) {
	e := NewJORMexplorer(path, command, coin)
	//j.ServicesSRV(*service, *port, *coin)
	e.ExploreCoin()
	ticker := time.NewTicker(23 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				e.ExploreCoin()
				fmt.Println("JORM explorer wooikos")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	go log.Fatal().Err(e.WWW.ListenAndServe())

}
