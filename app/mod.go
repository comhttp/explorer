package app

import (
	"net/http"
	"time"

	"github.com/comhttp/jorm/mod/nodes"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/jdb"
	"github.com/comhttp/jorm/pkg/utl"
)

type JORMexplorer struct {
	Coin       string
	BitNodes   nodes.BitNodes
	Status     *BlockchainStatus
	EQ         *ExplorerQueries
	config     cfg.Config
	jdbServers map[string]string
	WWW        *http.Server
	command    string
}

func NewJORMexplorer(path, command, coin string) *JORMexplorer {
	e := new(JORMexplorer)

	// bitNodes := nodes.BitNodes{}
	//if err := app.CFG.Read("nodes", coin, &bitNodes); err != nil {
	//	log.Print("Error", err)
	//}
	e.Coin = coin
	// e.BitNodes = bitNodes
	e.config.Path = path

	c, _ := cfg.NewCFG(e.config.Path, nil)
	// e.config = cfg.Config{}
	err := c.Read("conf", "conf", &e.config)
	utl.ErrorLog(err)

	e.jdbServers = make(map[string]string)
	err = c.Read("conf", "jdbs", &e.jdbServers)
	utl.ErrorLog(err)
	err = c.Read("nodes", coin, &e.BitNodes)
	utl.ErrorLog(err)

	e.Queries(coin, "")

	e.WWW = &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	e.WWW.Addr = ":" + e.config.Port[coin]
	//ttt := j.JDBS.B["coins"].ReadAllPerPages("coin", 10, 1)
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

//type Explorer struct {
//	Status map[string]*BlockchainStatus `json:"status"`
//}
type ExplorerQueries struct {
	info   *jdb.JDB
	blocks *jdb.JDB
	txs    *jdb.JDB
	addrs  *jdb.JDB
	col    string
}

type Explorer struct {
	Coin   string
	Status *BlockchainStatus
}
type Explorers map[string]*Explorer

type BlockchainStatus struct {
	Blocks    int `json:"blocks"`
	Txs       int `json:"txs"`
	Addresses int `json:"addresses"`
}

func (eq *ExplorerQueries) NewExplorer(coin string) *Explorer {

	return nil
}
func (e *JORMexplorer) Queries(coin, col string) {
	e.EQ = &ExplorerQueries{
		// &BlockchainStatus{},
		col: col,
	}
	info, err := jdb.NewJDB(e.jdbServers[coin])
	utl.ErrorLog(err)
	e.EQ.info = info
	blocks, err := jdb.NewJDB(e.jdbServers[coin+"blocks"])
	utl.ErrorLog(err)
	e.EQ.blocks = blocks
	if e.command != "onlyblocks" {
		txs, err := jdb.NewJDB(e.jdbServers[coin+"txs"])
		utl.ErrorLog(err)
		e.EQ.txs = txs
		addrs, err := jdb.NewJDB(e.jdbServers[coin+"addrs"])
		utl.ErrorLog(err)
		e.EQ.addrs = addrs
	}
	return
}
