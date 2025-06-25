package coin

import "testing"

func TestBlockValidate(t *testing.T) {
	tx := Transaction{Version: 1}
	blk := Block{Header: BlockHeader{Version: 1}, Transactions: []Transaction{tx}}
	blk.Header.MerkleRoot = blk.BuildMerkleRoot()
	if err := blk.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestBlockValidateMerkleError(t *testing.T) {
	tx := Transaction{Version: 1}
	blk := Block{Header: BlockHeader{Version: 1}, Transactions: []Transaction{tx}}
	blk.Header.MerkleRoot = "bad"
	if err := blk.Validate(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestBlockValidateDuplicateTx(t *testing.T) {
	tx := Transaction{Version: 1}
	blk := Block{
		Header:       BlockHeader{Version: 1},
		Transactions: []Transaction{tx, tx},
	}
	blk.Header.MerkleRoot = blk.BuildMerkleRoot()
	if err := blk.Validate(); err == nil {
		t.Fatalf("expected error")
	}
}
