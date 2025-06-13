package coin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// PointOut references a previous transaction output.
type PointOut struct {
	Hash  string `json:"hash"`
	Index uint32 `json:"index"`
}

// TxIn represents a transaction input.
type TxIn struct {
	PreviousOut PointOut `json:"prev_out"`
	ScriptSig   []byte   `json:"script_sig"`
	Sequence    uint32   `json:"sequence"`
}

// TxOut represents a transaction output.
type TxOut struct {
	Value        int64  `json:"value"`
	ScriptPubKey []byte `json:"script_pub_key"`
}

// Transaction is a simplified transaction structure.
type Transaction struct {
	Version  uint32  `json:"version"`
	Inputs   []TxIn  `json:"vin"`
	Outputs  []TxOut `json:"vout"`
	LockTime uint32  `json:"lock_time"`
}

// Hash returns the sha256 hash of the serialized transaction.
func (tx Transaction) Hash() string {
	data, _ := json.Marshal(tx)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
