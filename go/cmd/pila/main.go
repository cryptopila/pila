package main

import (
	"flag"
	"fmt"
	"log"

	"pila/pkg/coin"
	"pila/pkg/database"
)

func main() {
	list := flag.Bool("list", false, "list blocks")
	dbPath := flag.String("db", "./db", "database path")
	flag.Parse()

	db, err := database.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *list {
		blocks, err := db.ListBlocks()
		if err != nil {
			log.Fatal(err)
		}
		for _, b := range blocks {
			fmt.Printf("block %s with %d txs\n", b.Header.Hash()[:8], len(b.Transactions))
		}
		return
	}

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
			Version:   1,
			PrevHash:  "0",
			Timestamp: 0,
			Bits:      0,
			Nonce:     0,
		},
		Transactions: []coin.Transaction{tx},
	}
	blk.Header.MerkleRoot = blk.BuildMerkleRoot()

	if err := db.PutBlock(blk); err != nil {
		log.Fatal(err)
	}
	out, err := db.GetBlock(blk.Header.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("pila go stub running - loaded block %s with %d txs\n",
		out.Header.Hash()[:8], len(out.Transactions))
}
