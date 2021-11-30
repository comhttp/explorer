package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/jdb"
	"github.com/comhttp/jorm/pkg/utl"
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
func NewJORMexplorer(path, command, coin string) *JORMexplorer {
	e := new(JORMexplorer)
	e.Coin = coin
	e.config.Path = path
	e.command = command
	c, _ := cfg.NewCFG(e.config.Path, nil)

	err := c.Read("conf", "conf", &e.config)
	utl.ErrorLog(err)

	e.jdbServers = make(map[string]string)
	err = c.Read("conf", "jdbs", &e.jdbServers)
	utl.ErrorLog(err)

	err = c.Read("nodes", coin, &e.BitNodes)
	utl.ErrorLog(err)
	e.EQ = e.Queries()
	e.WWW = &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	e.WWW.Addr = ":" + e.config.Port[coin]
	status, err := e.GetStatus(coin)
	utl.ErrorLog(err)
	e.Status = status
	if e.Status == nil {
		e.Status = &BlockchainStatus{
			Blocks:    0,
			Txs:       0,
			Addresses: 0,
		}
	}
	return e
}
func (e *JORMexplorer) Queries() *ExplorerQueries {
	eq := &ExplorerQueries{
		// &BlockchainStatus{},
		// col: col,
	}

	info, err := jdb.NewJDB(e.jdbServers[e.Coin])
	utl.ErrorLog(err)
	eq.info = info
	blocks, err := jdb.NewJDB(e.jdbServers[e.Coin+"blocks"])
	utl.ErrorLog(err)
	eq.blocks = blocks
	if e.command != "onlyblocks" {
		txs, err := jdb.NewJDB(e.jdbServers[e.Coin+"txs"])
		utl.ErrorLog(err)
		eq.txs = txs
		addrs, err := jdb.NewJDB(e.jdbServers[e.Coin+"addrs"])
		utl.ErrorLog(err)
		eq.addrs = addrs
	}
	return eq
}
