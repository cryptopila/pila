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

func TestDBListBlocks(t *testing.T) {
	dir := t.TempDir()
	db, err := Open(dir)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	mkblk := func(v int) coin.Block {
		tx := coin.Transaction{Version: uint32(v)}
		b := coin.Block{Header: coin.BlockHeader{Version: 1}, Transactions: []coin.Transaction{tx}}
		b.Header.MerkleRoot = b.BuildMerkleRoot()
		return b
	}

	b1 := mkblk(1)
	b2 := mkblk(2)
	if err := db.PutBlock(b1); err != nil {
		t.Fatalf("put1: %v", err)
	}
	if err := db.PutBlock(b2); err != nil {
		t.Fatalf("put2: %v", err)
	}

	blocks, err := db.ListBlocks()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	seen := make(map[string]bool)
	for _, b := range blocks {
		seen[b.Header.Hash()] = true
	}
	if !seen[b1.Header.Hash()] || !seen[b2.Header.Hash()] {
		t.Fatalf("missing blocks")
	}
}
