package coin

import (
	"encoding/hex"
	"encoding/json"
)

// BlockHeader mirrors the basic Bitcoin block header structure.
type BlockHeader struct {
	Version    uint32 `json:"version"`
	PrevHash   string `json:"prev_hash"`
	MerkleRoot string `json:"merkle_root"`
	Timestamp  uint32 `json:"timestamp"`
	Bits       uint32 `json:"bits"`
	Nonce      uint32 `json:"nonce"`
}

// Hash returns the sha256d of the encoded header as a hex string.
func (h BlockHeader) Hash() string {
	data, _ := json.Marshal(h)
	sum := DoubleSHA256(data)
	return hex.EncodeToString(sum[:])
}

// Block groups a header with a list of transactions.
type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"tx"`
}

// BuildMerkleRoot calculates the merkle root of the block transactions.
func (b Block) BuildMerkleRoot() string {
	if len(b.Transactions) == 0 {
		return ""
	}

	var layer [][]byte
	for _, tx := range b.Transactions {
		layer = append(layer, tx.hashBytes())
	}

	for len(layer) > 1 {
		if len(layer)%2 != 0 {
			layer = append(layer, layer[len(layer)-1])
		}
		var next [][]byte
		for i := 0; i < len(layer); i += 2 {
			combined := append(append([]byte{}, layer[i]...), layer[i+1]...)
			h := DoubleSHA256(combined)
			next = append(next, h[:])
		}
		layer = next
	}

	return hex.EncodeToString(layer[0])
}
