package coin

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

// bytes serializes the header in the standard Bitcoin format.
func (h BlockHeader) bytes() []byte {
	buf := make([]byte, 80)
	binary.LittleEndian.PutUint32(buf[0:4], h.Version)
	prev, _ := hex.DecodeString(h.PrevHash)
	copy(buf[4:36], prev)
	root, _ := hex.DecodeString(h.MerkleRoot)
	copy(buf[36:68], root)
	binary.LittleEndian.PutUint32(buf[68:72], h.Timestamp)
	binary.LittleEndian.PutUint32(buf[72:76], h.Bits)
	binary.LittleEndian.PutUint32(buf[76:80], h.Nonce)
	return buf
}

// Hash returns the block header hash as a hex string.  For historical
// compatibility with the original C++ code, blocks with version < 5 use a
// Whirlpool-based hash while newer versions use Blake-256.
func (h BlockHeader) Hash() string {
	var digest [32]byte
	header := h.bytes()
	if h.Version < 5 {
		digest = WhirlpoolX(header)
	} else {
		digest = Blake256EightRound(header)
	}
	return hex.EncodeToString(digest[:])
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

// Validate performs basic sanity checks on the block.
// It currently verifies the merkle root matches the transactions.
func (b Block) Validate() error {
	if len(b.Transactions) == 0 {
		return fmt.Errorf("no transactions")
	}
	if b.Header.MerkleRoot != b.BuildMerkleRoot() {
		return fmt.Errorf("invalid merkle root")
	}
	seen := make(map[string]struct{})
	for _, tx := range b.Transactions {
		h := tx.Hash()
		if _, ok := seen[h]; ok {
			return fmt.Errorf("duplicate transaction %s", h)
		}
		seen[h] = struct{}{}
	}
	return nil
}
