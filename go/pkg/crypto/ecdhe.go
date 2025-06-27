package crypto

import (
	"github.com/btcsuite/btcd/btcec/v2"
)

// ECDHE implements a minimal secp256k1 Elliptic Curve Diffie-Hellman exchange.
type ECDHE struct {
	priv *btcec.PrivateKey
	pub  *btcec.PublicKey
}

// NewECDHE generates a fresh ECDHE instance using secp256k1.
func NewECDHE() (*ECDHE, error) {
	priv, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	return &ECDHE{priv: priv, pub: priv.PubKey()}, nil
}

// Public returns the compressed public key for this party.
func (e *ECDHE) Public() []byte {
	return e.pub.SerializeCompressed()
}

// Derive computes the shared secret given a peer's compressed public key.
func (e *ECDHE) Derive(peer []byte) ([]byte, error) {
	pub, err := btcec.ParsePubKey(peer)
	if err != nil {
		return nil, err
	}
	secret := btcec.GenerateSharedSecret(e.priv, pub)
	return secret, nil
}
