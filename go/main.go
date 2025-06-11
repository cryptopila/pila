package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/cryptopila/pila/pkg"
)

func main() {
	db, err := pkg.NewDB("./data")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx := pkg.Transaction{ID: []byte("tx1"), From: "a", To: "b", Amount: 10}
	block := pkg.Block{PrevHash: []byte("genesis"), Txns: []pkg.Transaction{tx}}
	hash := sha256.Sum256([]byte("block1"))
	block.Hash = hash[:]

	if err := db.PutBlock(&block); err != nil {
		panic(err)
	}
	b, err := db.GetBlock(block.Hash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("stored block with %d txns\n", len(b.Txns))
}
