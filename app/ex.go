package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/comhttp/jorm/mod/nodes"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/jdb"
	"github.com/comhttp/jorm/pkg/strapi"
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
	e := &JORMexplorer{}
	jdbServers := make(map[string]string)
	e.Coin = coin
	e.config.Path = path
	e.command = command
	c, _ := cfg.NewCFG(e.config.Path, nil)

	err := c.Read("conf", "conf", &e.config)
	utl.ErrorLog(err)

	// e.jdbServers = make(map[string]string)
	// err = c.Read("conf", "jdbs", &e.jdbServers)
	// utl.ErrorLog(err)

	// err = c.Read("nodes", coin, &e.BitNodes)
	// utl.ErrorLog(err)

	e.okno = strapi.New(e.config.Strapi)

	bitnodes := e.okno.GetAll("nodes", "bitnode=true&")

	for _, bitnode := range bitnodes {
		if bitnode["coin"].(map[string]interface{})["slug"] == coin {
			e.BitNodes = append(e.BitNodes, nodes.BitNode{
				IP:   bitnode["ip"].(string),
				Port: int64(bitnode["port"].(float64)),
			})
		}
	}
	jdbs := e.okno.GetAll("services", "type=jdb&")
	for _, jdb := range jdbs {
		jdbServers[jdb["slug"].(string)] = jdb["server"].(map[string]interface{})["tailscale"].(string) + ":" + fmt.Sprint(jdb["port"].(float64))
	}
	e.ExJDBs = InitExplorerJDBs(jdbServers, e.command, e.Coin)
	e.WWW = &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	e.WWW.Addr = ":" + e.config.Port[coin]
	status, err := e.ExJDBs.GetStatus(coin)
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
func InitExplorerJDBs(jdbServers map[string]string, command, coin string) *ExplorerJDBs {
	eq := &ExplorerJDBs{
		// &BlockchainStatus{},
		// col: col,
	}
	info, err := jdb.NewJDB(jdbServers["jdb"+coin])
	utl.ErrorLog(err)
	eq.info = info
	blocks, err := jdb.NewJDB(jdbServers["jdb"+coin+"blocks"])
	utl.ErrorLog(err)
	eq.blocks = blocks
	if command != "onlyblocks" {
		txs, err := jdb.NewJDB(jdbServers["jdb"+coin+"txs"])
		utl.ErrorLog(err)
		eq.txs = txs
		addrs, err := jdb.NewJDB(jdbServers["jdb"+coin+"addrs"])
		utl.ErrorLog(err)
		eq.addrs = addrs
	}
	return eq
}
