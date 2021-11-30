package app

import (
	"net/http"

	"github.com/comhttp/jorm/mod/nodes"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/jdb"
	"github.com/comhttp/jorm/pkg/strapi"
)

type JORMexplorer struct {
	Coin     string
	BitNodes nodes.BitNodes
	Status   *BlockchainStatus
	ExJDBs   *ExplorerJDBs
	config   cfg.Config
	WWW      *http.Server
	command  string
	okno     strapi.StrapiRestClient
}

type ExplorerJDBs struct {
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
