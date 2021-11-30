package app

import (
	"net/http"

	"github.com/comhttp/jorm/mod/nodes"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/jdb"
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

type ExplorerQueries struct {
	info   *jdb.JDB
	blocks *jdb.JDB
	txs    *jdb.JDB
	addrs  *jdb.JDB
	col    string
}

type BlockchainStatus struct {
	Blocks    int `json:"blocks"`
	Txs       int `json:"txs"`
	Addresses int `json:"addresses"`
}
