package crypto

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

// PBKDF2SHA256 derives a key from the given password and salt using
// the PBKDF2 algorithm with HMAC-SHA256.
func PBKDF2SHA256(password, salt []byte, iter int, keyLen int) []byte {
	return pbkdf2.Key(password, salt, iter, keyLen, sha256.New)
}
