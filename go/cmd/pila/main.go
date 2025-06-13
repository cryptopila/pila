package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"

	"pila/pkg/coin"
)

func main() {
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx := coin.Transaction{
		Version: 1,
		Inputs: []coin.TxIn{{
			PreviousOut: coin.PointOut{Hash: "prev", Index: 0},
			ScriptSig:   []byte("sig"),
			Sequence:    0xffffffff,
		}},
		Outputs: []coin.TxOut{{
			Value:        50,
			ScriptPubKey: []byte("pub"),
		}},
		LockTime: 0,
	}
	blk := coin.Block{
		Header: coin.BlockHeader{
			Version:    1,
			PrevHash:   "0",
			MerkleRoot: tx.Hash(),
			Timestamp:  0,
			Bits:       0,
			Nonce:      0,
		},
		Transactions: []coin.Transaction{tx},
	}

	data, err := json.Marshal(blk)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Put([]byte("block:1"), data, nil); err != nil {
		log.Fatal(err)
	}

	raw, err := db.Get([]byte("block:1"), nil)
	if err != nil {
		log.Fatal(err)
	}

	var out coin.Block
	if err := json.Unmarshal(raw, &out); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pila go stub running - loaded block with %d txs\n", len(out.Transactions))
}
