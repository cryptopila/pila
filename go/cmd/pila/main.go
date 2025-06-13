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

	blk := coin.Block{
		Hash:         "demo",
		Transactions: []coin.Transaction{{ID: "tx1"}},
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

	fmt.Printf("pila go stub running - loaded block %s with %d txs\n", out.Hash, len(out.Transactions))
}
