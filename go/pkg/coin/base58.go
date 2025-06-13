package coin

import (
	"crypto/sha256"
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

// Base58 represents a basic base58 encoded value with a version byte.
type Base58 struct {
	Version uint8
	Data    []byte
}

// SetData initializes the object from a version and data buffer.
func (b *Base58) SetData(version uint8, buf []byte) {
	b.Version = version
	b.Data = append([]byte(nil), buf...)
}

// SetString decodes a base58 string with checksum.
func checksum(data []byte) [4]byte {
	h := sha256.Sum256(data)
	h2 := sha256.Sum256(h[:])
	var out [4]byte
	copy(out[:], h2[:4])
	return out
}

func decodeBase58Check(str string) ([]byte, error) {
	decoded := base58.Decode(str)
	if len(decoded) < 4 {
		return nil, errors.New("invalid base58 format")
	}
	data := decoded[:len(decoded)-4]
	got := decoded[len(decoded)-4:]
	cksum := checksum(data)
	if !equalBytes(got, cksum[:]) {
		return nil, errors.New("checksum error")
	}
	return data, nil
}

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (b *Base58) SetString(val string) bool {
	tmp, err := decodeBase58Check(val)
	if err != nil || len(tmp) == 0 {
		b.Data = nil
		b.Version = 0
		return false
	}
	b.Version = tmp[0]
	b.Data = append([]byte(nil), tmp[1:]...)
	return true
}

// ToString encodes the object as base58 with checksum.
func encodeBase58Check(data []byte) string {
	cksum := checksum(data)
	buf := append(append([]byte{}, data...), cksum[:]...)
	return base58.Encode(buf)
}

func (b Base58) ToString(includeVersion bool) string {
	var v []byte
	if includeVersion {
		v = append([]byte{b.Version}, b.Data...)
	} else {
		v = b.Data
	}
	return encodeBase58Check(v)
}

// CompareTo compares two base58 values.
func (b Base58) CompareTo(o Base58) int {
	if b.Version < o.Version {
		return -1
	}
	if b.Version > o.Version {
		return 1
	}
	lb := len(b.Data)
	lo := len(o.Data)
	for i := 0; i < lb && i < lo; i++ {
		if b.Data[i] < o.Data[i] {
			return -1
		}
		if b.Data[i] > o.Data[i] {
			return 1
		}
	}
	switch {
	case lb < lo:
		return -1
	case lb > lo:
		return 1
	default:
		return 0
	}
}

func (b Base58) VersionByte() uint8 { return b.Version }
func (b Base58) Bytes() []byte      { return b.Data }
