package database

import (
	"testing"

	"pila/pkg/coin"
)

func TestDBPutGetBlock(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(dir)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	tx := coin.Transaction{Version: 1}
	blk := coin.Block{Header: coin.BlockHeader{Version: 1}, Transactions: []coin.Transaction{tx}}
	blk.Header.MerkleRoot = blk.BuildMerkleRoot()
	if err := db.PutBlock(blk); err != nil {
		t.Fatalf("put: %v", err)
	}
	out, err := db.GetBlock(blk.Header.Hash())
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if out.Header.Hash() != blk.Header.Hash() {
		t.Fatalf("mismatch: %s != %s", out.Header.Hash(), blk.Header.Hash())
	}
}

func TestDBPutBlockInvalid(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(dir)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	blk := coin.Block{Header: coin.BlockHeader{Version: 1}}
	if err := db.PutBlock(blk); err == nil {
		t.Fatalf("expected error")
	}
}
