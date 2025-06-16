package coin

import (
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

// Hash160 computes a ripemd160(sha256(data)) hash, similar to Bitcoin's HASH160.
func Hash160(data []byte) []byte {
	h := sha256.Sum256(data)
	r := ripemd160.New()
	r.Write(h[:])
	return r.Sum(nil)
}

// DoubleSHA256 computes sha256d(data) which is sha256(sha256(data)).
func DoubleSHA256(data []byte) [32]byte {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}
