package exp

import (
	"github.com/comhttp/jorm/pkg/cfg"
	"github.com/comhttp/jorm/pkg/utl"
	"github.com/comhttp/node/nodes"
	"strconv"
)

// GetExplorer updates the data from blockchain of a coin in the database
func ExploreCoin(e *JORMexplorer, coin string) {
	var b []string
	log.Println("Coin is BitNode:", coin)
	if e.BitNodes != nil {
		b = append(b, coin)
		for _, bitnode := range e.BitNodes {
			bitnode.Jrc = utl.NewClient(cfg.C.RPC.Username, cfg.C.RPC.Password, bitnode.IP, bitnode.Port)
			log.Println("Get Coin Blockchain:", coin)
			go e.GetCoinBlockchain(&bitnode, coin)
		}
	}
}

// GetExplorer returns the full set of information about a block
func (e *JORMexplorer) GetCoinBlockchain(b *nodes.BitNode, c string) {
	if b.Jrc != nil {
		blockCount := b.APIGetBlockCount()
		log.Println("Block Count from the chain: ", blockCount)
		log.Println("Status :  :::"+c+" - - ", e.status.Blocks)
		if blockCount >= e.status.Blocks {
			e.blocks(b, blockCount, c)
		}
	}
}

func (e *JORMexplorer) blocks(b *nodes.BitNode, bc int, c string) {
	for {
		blockRaw := b.APIGetBlockByHeight(e.status.Blocks)
		if blockRaw != nil && blockRaw != "" {
			blockHash := blockRaw.(map[string]interface{})["hash"].(string)
			e.JDB.Write(c, "block_"+strconv.Itoa(e.status.Blocks), blockHash)
			e.JDB.Write(c, "block_"+blockHash, blockRaw)
			block := (blockRaw).(map[string]interface{})
			if e.status.Blocks != 0 {
				for _, t := range (block["tx"]).([]interface{}) {
					e.tx(b, c, t.(string))
				}
			}
			bl := blockRaw.(map[string]interface{})
			e.status.Blocks = int(bl["height"].(float64))
			log.Println("Write "+c+" block: "+strconv.Itoa(e.status.Blocks)+" - ", blockHash)
			e.JDB.Write(c, "status", e.status)
		} else {
			break
		}
		if bc != 0 {
			e.status.Blocks++
		}
		log.Println("StatusBlocks   "+c, e.status.Blocks)
	}
}

func (e *JORMexplorer) tx(b *nodes.BitNode, c, txid string) {
	txRaw := b.APIGetTx(txid)
	e.status.Txs++
	e.JDB.Write(c, "tx_"+txid, txRaw)
	if txRaw != nil {
		tx := (txRaw).(map[string]interface{})
		if tx["vout"] != nil {
			for _, nRaw := range tx["vout"].([]interface{}) {
				if nRaw.(map[string]interface{})["scriptPubKey"] != nil {
					scriptPubKey := nRaw.(map[string]interface{})["scriptPubKey"].(map[string]interface{})
					if scriptPubKey["addresses"] != nil {
						for _, address := range (scriptPubKey["addresses"]).([]interface{}) {
							e.addr(c, address.(string))
						}
					}
				}
			}
		}
	}
	log.Println("Write "+c+" transaction: ", txid)
	return
}

func (e *JORMexplorer) addr(c, address string) {
	e.status.Addresses++
	e.JDB.Write(c, "addr_"+address, address)
	log.Println("Write "+c+" address: ", address)
	return
}

//func (e *explorer.Explorer) status(n *nodes.BitNode) {
//
//	fmt.Println("Mempool: ", n.GetRawMemPool())
//	fmt.Println("MiningInfo: ", n.GetMiningInfo())
//	fmt.Println("NetworkInfo: ", n.GetNetworkInfo())
//	fmt.Println("Info: ", n.GetInfo())
//	fmt.Println("PeerInfo: ", n.GetPeerInfo())
//
//
//	//n.addNode(ip string)
//	//n.GetAddNodeInfo(ip string)
//}
//
//
