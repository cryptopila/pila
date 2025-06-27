package crypto

import (
	"encoding/hex"
	"testing"
)

func TestPBKDF2SHA256(t *testing.T) {
	password := []byte("password")
	salt := []byte("salt")
	dk := PBKDF2SHA256(password, salt, 4096, 32)
	expected := "c5e478d59288c841aa530db6845c4c8d962893a001ce4e11a4963873aa98134a"
	if hex.EncodeToString(dk) != expected {
		t.Fatalf("unexpected result: %x", dk)
	}
}
