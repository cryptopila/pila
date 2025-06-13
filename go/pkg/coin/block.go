package coin

// BlockHeader mirrors the basic Bitcoin block header structure.
type BlockHeader struct {
	Version    uint32 `json:"version"`
	PrevHash   string `json:"prev_hash"`
	MerkleRoot string `json:"merkle_root"`
	Timestamp  uint32 `json:"timestamp"`
	Bits       uint32 `json:"bits"`
	Nonce      uint32 `json:"nonce"`
}

// Block groups a header with a list of transactions.
type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"tx"`
}
