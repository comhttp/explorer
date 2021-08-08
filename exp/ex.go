package exp

import (
	"fmt"
	"github.com/comhttp/jorm/mod/explorers"
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/jdb"
	"github.com/comhttp/jorm/pkg/utl"
	"github.com/comhttp/node/nodes"
	"github.com/sirupsen/logrus"
)

//type Explorer struct {
//	Status map[string]*BlockchainStatus `json:"status"`
//}

type JORMexplorer struct {
	Coin     string
	BitNodes nodes.BitNodes
	status   *explorers.BlockchainStatus
	JDB      *jdb.JDB
}

func NewJORMexplorer(loglevel, coin string) *JORMexplorer {
	log.SetLevel(parseLogLevel(loglevel))
	err := cfg.CFG.Read("conf", "conf", &cfg.C)
	utl.ErrorLog(err)
	bitNodes := nodes.BitNodes{}
	if err := cfg.CFG.Read("nodes", coin, &bitNodes); err != nil {
		fmt.Println("Error", err)
	}
	e := &JORMexplorer{
		Coin:     coin,
		BitNodes: bitNodes,
		JDB:      jdb.NewJDB(cfg.C.JDBservers),
	}
	err = e.JDB.Read(coin, "status", &e.status)
	if e.status == nil {
		e.status = &explorers.BlockchainStatus{
			Blocks:    0,
			Txs:       0,
			Addresses: 0,
		}
	}
	utl.ErrorLog(err)
	return e
}

var log = logrus.New()

func wrapLogger(module string) logrus.FieldLogger {
	return log.WithField("module", module)
}

func parseLogLevel(level string) logrus.Level {
	switch level {
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info", "notice":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
