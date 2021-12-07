package app

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/comhttp/jorm/pkg/utl"
)

//
//type Explorers struct {
//	Status map[string]*BlockchainStatus `json:"status"`
//}

//func GetExplorers(j *jdb.JDBS, nodeCoins []string) {
//status: &BlockchainStatus{},
//for _, coin := range nodeCoins {
//	s := &BlockchainStatus{}
//err := j.coin[coin].Read(coin, "status", &s)
//utl.ErrorLog(err)
//j.Explorers[coin].status = s
//}
//return
//}

func (ej *ExplorerJDBs) GetExplorer(coin string) *BlockchainStatus {
	s := &BlockchainStatus{}
	err := ej.info.Read(coin, "status", &s)
	utl.ErrorLog(err)
	return s
}

func (ej *ExplorerJDBs) GetStatus(status *BlockchainStatus, coin string) (*BlockchainStatus, error) {
	err := ej.info.Read(coin, "status", &status)
	utl.ErrorLog(err)
	fmt.Println("ej.status", status)
	return status, err
}

func (ej *ExplorerJDBs) GetLastBlock(status *BlockchainStatus, coin string) int {
	status, err := ej.GetStatus(status, coin)
	utl.ErrorLog(err)
	return status.Blocks
}

func (ej *ExplorerJDBs) GetBlock(coin, id string) map[string]interface{} {
	blockHash := ""
	block := make(map[string]interface{})
	_, err := strconv.Atoi(id)
	if err != nil {
		blockHash = id
	} else {
		blockHash = ""
		err = ej.blocks.Read("block", id, &blockHash)
	}
	err = ej.blocks.Read("block", blockHash, &block)
	utl.ErrorLog(err)
	return block
}

func (ej *ExplorerJDBs) GetBlocks(coin string, per, page int) (blocks []map[string]interface{}) {
	s := &BlockchainStatus{}
	err := ej.info.Read(coin, "status", &s)
	utl.ErrorLog(err)
	blockCount := s.Blocks
	//app.log.Print("blockCount", blockCount)
	startBlock := blockCount - per*page
	minusBlockStart := int(startBlock + per)
	for ibh := minusBlockStart; ibh >= startBlock; {
		blocks = append(blocks, ej.GetBlockShort(coin, strconv.Itoa(ibh)))
		ibh--
	}
	sort.SliceStable(blocks, func(i, j int) bool {
		return int(blocks[i]["height"].(int64)) > int(blocks[j]["height"].(int64))
	})
	return blocks
}
func (ej *ExplorerJDBs) GetBlockShort(coin, blockhash string) map[string]interface{} {
	b := ej.GetBlock(coin, blockhash)
	block := make(map[string]interface{})
	if b["bits"] != nil {
		block["bits"] = b["bits"].(string)
	}
	if b["confirmations"] != nil {
		block["confirmations"] = int64(b["confirmations"].(float64))
	}
	if b["difficulty"] != nil {
		block["difficulty"] = b["difficulty"].(float64)
	}
	if b["hash"] != nil {
		block["hash"] = b["hash"].(string)
	}
	if b["height"] != nil {
		block["height"] = int64(b["height"].(float64))
	}
	if b["tx"] != nil {
		var txsNumber int
		for _ = range b["tx"].([]interface{}) {
			txsNumber++
		}
		block["txs"] = txsNumber
	}
	if b["size"] != nil {
		block["size"] = int64(b["size"].(float64))
	}
	if b["time"] != nil {
		unixTimeUTC := time.Unix(int64(b["time"].(float64)), 0)
		block["time"] = unixTimeUTC.Format(time.RFC850)
		block["timeutc"] = unixTimeUTC.Format(time.RFC3339)
	}
	return block
}

func (ej *ExplorerJDBs) GetTx(coin, id string) map[string]interface{} {
	tx := make(map[string]interface{})
	err := ej.txs.Read("tx", id, &tx)
	utl.ErrorLog(err)
	return tx
}
func (ej *ExplorerJDBs) GetAddr(coin, id string) map[string]interface{} {
	addr := make(map[string]interface{})
	err := ej.addrs.Read("addr", id, &addr)
	utl.ErrorLog(err)
	return addr
}

func (ej *ExplorerJDBs) GetMemPool(coin string) []string {
	mempool := []string{}
	err := ej.info.Read(coin, "mempool", &mempool)
	utl.ErrorLog(err)
	return mempool
}

func (ej *ExplorerJDBs) GetMiningInfo(coin string) map[string]interface{} {
	mininginfo := make(map[string]interface{})
	err := ej.info.Read(coin, "mining", &mininginfo)
	utl.ErrorLog(err)
	return mininginfo
}

func (ej *ExplorerJDBs) GetInfo(coin string) map[string]interface{} {
	info := make(map[string]interface{})
	err := ej.info.Read(coin, "info", &info)
	utl.ErrorLog(err)
	return info
}

func (ej *ExplorerJDBs) GetNetworkInfo(coin string) map[string]interface{} {
	network := make(map[string]interface{})
	err := ej.info.Read(coin, "network", &network)
	utl.ErrorLog(err)
	return network
}

func (ej *ExplorerJDBs) GetPeers(coin string) []interface{} {
	peers := new([]interface{})
	err := ej.info.Read(coin, "peers", &peers)
	utl.ErrorLog(err)
	return *peers
}
