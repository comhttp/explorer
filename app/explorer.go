package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/comhttp/jorm/mod/nodes"
	"github.com/comhttp/jorm/pkg/utl"
	"github.com/rs/zerolog/log"
)

// GetExplorer updates the data from blockchain of a coin in the database
func (e *JORMexplorer) ExploreCoin() {
	log.Print("Coin is BitNode:", e.Coin)
	if e.BitNodes != nil {
		for _, bitnode := range e.BitNodes {
			log.Print("Bitnode: ", bitnode)
			bitnode.Jrc = utl.NewClient(e.config.RPC.Username, e.config.RPC.Password, bitnode.IP, bitnode.Port)
			e.EQ.status.Write(e.Coin, "info", bitnode.APIGetInfo())
			e.EQ.status.Write(e.Coin, "peers", bitnode.APIGetPeerInfo())
			e.EQ.status.Write(e.Coin, "mempool", bitnode.APIGetRawMemPool())
			e.EQ.status.Write(e.Coin, "mining", bitnode.APIGetMiningInfo())
			e.EQ.status.Write(e.Coin, "network", bitnode.APIGetNetworkInfo())

			log.Print("Get Coin Blockchain:", e.Coin)
			e.blockchain(&bitnode, e.Coin)
		}
	}
}

// GetExplorer returns the full set of information about a block
func (e *JORMexplorer) blockchain(bn *nodes.BitNode, coin string) {
	if bn.Jrc != nil {
		blockCount := bn.APIGetBlockCount()

		fmt.Println("-------------------------------------")
		fmt.Println("blockCount", blockCount)
		fmt.Println("e.Status.Blocks", e.Status.Blocks)
		fmt.Println("-------------------------------------")
		log.Print("Block Count from the chain: ", blockCount)
		log.Print("Status ::: "+coin+" ::: ", e.Status.Blocks)
		if blockCount >= e.Status.Blocks {
			// e.block(bn, coin)
			e.blocks(bn, blockCount, coin)
		}
	}
}

func (e *JORMexplorer) block(b *nodes.BitNode, coin string) {
	blockRaw := b.APIGetBlockByHeight(e.Status.Blocks - 1)
	if blockRaw != nil && blockRaw != "" {
		blockHash := blockRaw.(map[string]interface{})["hash"].(string)
		e.EQ.blocks.Write("block", strconv.Itoa(e.Status.Blocks), blockHash)
		e.EQ.blocks.Write("block", blockHash, blockRaw)
		block := (blockRaw).(map[string]interface{})
		if e.command != "onlyblocks" {
			go e.txs(b, (block["tx"]).([]interface{}), coin)
		}

		bl := blockRaw.(map[string]interface{})
		e.Status.Blocks = int(bl["height"].(float64))
		log.Print("Write "+coin+" block: "+strconv.Itoa(e.Status.Blocks)+" - ", blockHash)
		e.EQ.status.Write(e.Coin, "status", e.Status)
	}
}

func (e *JORMexplorer) blocks(b *nodes.BitNode, bc int, coin string) {
	for {
		blockRaw := b.APIGetBlockByHeight(e.Status.Blocks)
		if blockRaw != nil && blockRaw != "" {
			blockHash := blockRaw.(map[string]interface{})["hash"].(string)
			e.EQ.blocks.Write("block", strconv.Itoa(e.Status.Blocks), blockHash)
			e.EQ.blocks.Write("block", blockHash, blockRaw)
			block := (blockRaw).(map[string]interface{})
			if e.Status.Blocks != 0 {
				if e.command != "onlyblocks" {
					go e.txs(b, (block["tx"]).([]interface{}), coin)
				}
			}
			bl := blockRaw.(map[string]interface{})
			e.Status.Blocks = int(bl["height"].(float64))
			log.Info().Msg("Write " + coin + " block: " + strconv.Itoa(e.Status.Blocks) + " - " + blockHash)
			e.EQ.status.Write(e.Coin, "status", e.Status)
		} else {
			break
		}
		if bc != 0 {
			e.Status.Blocks++
		}
		log.Print("StatusBlocks   "+coin+": ", e.Status.Blocks)
		time.Sleep(9 * time.Millisecond)
	}
}

func (e *JORMexplorer) tx(b *nodes.BitNode, coin, txid string) {
	txRaw := b.APIGetTx(txid)
	e.Status.Txs++
	e.EQ.txs.Write("tx", txid, txRaw)
	log.Info().Msg("Write " + coin + " transaction: " + txid)
	if txRaw != nil {
		tx := (txRaw).(map[string]interface{})
		// if tx["vin"] != nil {
		// 	// fmt.Println(": vin:", tx["vin"])
		// }

		// // fmt.Println(": txid:", tx["txid"])

		// for _, vout := range tx["vout"].([]interface{}) {
		// 	// 	// fmt.Println(": ttt:", ttt)

		// 	if vout.(map[string]interface{})["scriptPubKey"] != nil {
		// 		// n := ttt.(map[string]interface{})["n"].(float64)
		// 		// fmt.Println(": nnn:", n)

		// 		scriptPubKey := vout.(map[string]interface{})["scriptPubKey"].(map[string]interface{})

		// 		if scriptPubKey["addresses"] != nil {

		// 			addresses := scriptPubKey["addresses"].([]interface{})
		// 			fmt.Println(": addresse: ", addresses[0])
		// 			fmt.Println(": vout :", vout.(map[string]interface{})["value"])

		// 		}

		// 	}

		// }
		if tx["vout"] != nil {
			for _, nRaw := range tx["vout"].([]interface{}) {
				if nRaw.(map[string]interface{})["scriptPubKey"] != nil {
					scriptPubKey := nRaw.(map[string]interface{})["scriptPubKey"].(map[string]interface{})
					if scriptPubKey["addresses"] != nil {
						go e.addrs(b, (scriptPubKey["addresses"]).([]interface{}), nRaw.(map[string]interface{}), txid, coin)
					}
				}
			}
		}
	}
	return
}

func (e *JORMexplorer) addr(coin, address, txid string, value float64) {
	e.Status.Addresses++
	addr := make(map[string]interface{})
	err := e.EQ.addrs.Read("addr", address, &addr)
	utl.ErrorLog(err)
	addr["address"] = address
	if addr["value"] == nil {
		addr["value"] = 0.0
	}
	if addr["txs"] != nil {
		addr["txs"] = append(addr["txs"].([]interface{}), txid)
	} else {
		txs := []string{}
		txs = append(txs, txid)
		addr["txs"] = txs
	}
	addr["value"] = addr["value"].(float64) + value
	e.EQ.addrs.Write("addr", address, addr)
	log.Info().Msg("Write " + coin + " address: " + address)
	return
}

func (e *JORMexplorer) txs(b *nodes.BitNode, tx []interface{}, coin string) {
	for _, t := range tx {
		e.tx(b, coin, t.(string))
	}
}
func (e *JORMexplorer) addrs(b *nodes.BitNode, addresses []interface{}, nRaw map[string]interface{}, txid, coin string) {
	for _, address := range addresses {
		e.addr(coin, address.(string), txid, nRaw["value"].(float64))
	}
}
