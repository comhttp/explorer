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
	fmt.Println("-------------111111111111111111111------------------------")

	log.Print("Coin is BitNode:", e.Coin)
	if e.BitNodes != nil {
		for _, bitnode := range e.BitNodes {

			fmt.Println("-------------bbbbbbbbbbbbbbbbbbbbbbb------------------------")

			log.Print("Bitnode: ", bitnode)
			bitnode.Jrc = utl.NewClient(e.config.RPC.Username, e.config.RPC.Password, bitnode.IP, bitnode.Port)
			e.ExJDBs.info.Write(e.Coin, "info", bitnode.APIGetInfo())
			e.ExJDBs.info.Write(e.Coin, "peers", bitnode.APIGetPeerInfo())
			e.ExJDBs.info.Write(e.Coin, "mempool", bitnode.APIGetRawMemPool())
			e.ExJDBs.info.Write(e.Coin, "mining", bitnode.APIGetMiningInfo())
			e.ExJDBs.info.Write(e.Coin, "network", bitnode.APIGetNetworkInfo())

			log.Print("Get Coin Blockchain:", e.Coin)
			e.blockchain(&bitnode, e.Coin)
		}
	}
}

// GetExplorer returns the full set of information about a block
func (e *JORMexplorer) blockchain(bn *nodes.BitNode, coin string) {

	fmt.Println("-------------1212------------------------")

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
	saveBlockRaw := make(map[string]interface{})
	saveTxsRaw := make(map[string]interface{})
	saveAddrsRaw := make(map[string]interface{})
	if blockRaw != nil && blockRaw != "" {
		blockHash := blockRaw.(map[string]interface{})["hash"].(string)
		// blockSaveErr := e.EQ.blocks.Write("block", strconv.Itoa(e.Status.Blocks), blockHash)
		// blockSaveErr = e.EQ.blocks.Write("block", blockHash, blockRaw)

		saveBlockRaw["block_"+strconv.Itoa(e.Status.Blocks)] = blockHash
		saveBlockRaw["block_"+blockHash] = blockRaw

		err := e.ExJDBs.blocks.WriteAll(saveBlockRaw)
		err = e.ExJDBs.txs.WriteAll(saveTxsRaw)
		err = e.ExJDBs.addrs.WriteAll(saveAddrsRaw)

		block := (blockRaw).(map[string]interface{})
		if e.command != "onlyblocks" {
			e.txs(b, (block["tx"]).([]interface{}), coin, saveTxsRaw, saveAddrsRaw)
		}
		if err == nil {
			bl := blockRaw.(map[string]interface{})
			e.Status.Blocks = int(bl["height"].(float64))
			log.Print("Write "+coin+" block: "+strconv.Itoa(e.Status.Blocks)+" - ", blockHash)
			e.ExJDBs.info.Write(e.Coin, "status", e.Status)
		}

	}
}

func (e *JORMexplorer) blocks(b *nodes.BitNode, bc int, coin string) {
	for {
		blockRaw := b.APIGetBlockByHeight(e.Status.Blocks)
		saveBlockRaw := make(map[string]interface{})
		saveTxsRaw := make(map[string]interface{})
		saveAddrsRaw := make(map[string]interface{})
		if blockRaw != nil && blockRaw != "" {
			blockHash := blockRaw.(map[string]interface{})["hash"].(string)
			block := (blockRaw).(map[string]interface{})
			saveBlockRaw["block_"+strconv.Itoa(e.Status.Blocks)] = blockHash
			saveBlockRaw["block_"+blockHash] = blockRaw
			err := e.ExJDBs.blocks.WriteAll(saveBlockRaw)
			if e.Status.Blocks != 0 {
				if e.command != "onlyblocks" {
					e.txs(b, (block["tx"]).([]interface{}), coin, saveTxsRaw, saveAddrsRaw)
					err = e.ExJDBs.txs.WriteAll(saveTxsRaw)
					err = e.ExJDBs.addrs.WriteAll(saveAddrsRaw)
				}
			}
			if err == nil {
				bl := blockRaw.(map[string]interface{})
				e.Status.Blocks = int(bl["height"].(float64))
				log.Info().Msg("Write " + coin + " block: " + strconv.Itoa(e.Status.Blocks) + " - " + blockHash)
				e.ExJDBs.info.Write(e.Coin, "status", e.Status)
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
}

func (e *JORMexplorer) tx(b *nodes.BitNode, coin, txid string) (txRaw interface{}, saveAddrsRaw map[string]interface{}) {
	txRaw = b.APIGetTx(txid)
	e.Status.Txs++
	// e.EQ.txs.Write("tx", txid, txRaw)
	// log.Info().Msg("Write " + coin + " transaction: " + txid)
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
						// saveAddrsRaw = e.addrs(b, (scriptPubKey["addresses"]).([]interface{}), nRaw.(map[string]interface{}), txid, coin)
					}
				}
			}
		}
	}

	return txRaw, saveAddrsRaw
}

func (e *JORMexplorer) addr(coin, address, txid string, value float64) map[string]interface{} {
	e.Status.Addresses++
	addr := make(map[string]interface{})
	err := e.ExJDBs.addrs.Read("addr", address, &addr)
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
	// e.EQ.addrs.Write("addr", address, addr)
	// log.Info().Msg("Write " + coin + " address: " + address)
	return addr
}

func (e *JORMexplorer) txs(b *nodes.BitNode, txs []interface{}, coin string, saveTxsRaw, saveAddrsRaw map[string]interface{}) {
	// saveTxsRaw = make(map[string]interface{})
	// saveAddrsRaw = make(map[string]interface{})
	for _, t := range txs {
		txRaw, saveTxAddrsRaw := e.tx(b, coin, t.(string))
		for saveKey, saveData := range saveTxAddrsRaw {
			saveAddrsRaw[saveKey] = saveData
		}
		saveTxsRaw["tx_"+t.(string)] = txRaw
	}
	return
}
func (e *JORMexplorer) addrs(b *nodes.BitNode, addresses []interface{}, nRaw map[string]interface{}, txid, coin string) map[string]interface{} {
	saveAddrsRaw := make(map[string]interface{})
	for _, address := range addresses {
		saveAddrsRaw["addr_"+address.(string)] = e.addr(coin, address.(string), txid, nRaw["value"].(float64))
	}
	fmt.Println(":sa1111111111111111111veAddrsRawsaveAddrsRawsaveAddrsRawsaveAddrsRawsaveAddrsRaw: ", saveAddrsRaw)

	return saveAddrsRaw
}
