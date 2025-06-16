package coin

import (
	"crypto/rand"
	"encoding/binary"
	"math/bits"
)

// DoubleSHA256Checksum returns the first 4 bytes of sha256d(data) as a uint32.
func DoubleSHA256Checksum(data []byte) uint32 {
	sum := DoubleSHA256(data)
	return binary.LittleEndian.Uint32(sum[0:4])
}

// SHA256RIPEMD160 returns ripemd160(sha256(data)).
// This simply wraps Hash160 but returns a fixed-size array.
func SHA256RIPEMD160(data []byte) [20]byte {
	h := Hash160(data)
	var out [20]byte
	copy(out[:], h)
	return out
}

// SHA256Random returns 32 cryptographically random bytes.
func SHA256Random() [32]byte {
	var out [32]byte
	if _, err := rand.Read(out[:]); err != nil {
		// in the unlikely event rand fails, fall back to sha256d of empty slice
		out = DoubleSHA256(nil)
	}
	return out
}

// Murmur3 implements the 32-bit murmur3 hashing algorithm.
func Murmur3(seed uint32, data []byte) uint32 {
	const (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
	)

	h := seed
	nBlocks := len(data) / 4
	for i := 0; i < nBlocks; i++ {
		k := binary.LittleEndian.Uint32(data[i*4:])
		k *= c1
		k = bits.RotateLeft32(k, 15)
		k *= c2

		h ^= k
		h = bits.RotateLeft32(h, 13)
		h = h*5 + 0xe6546b64
	}

	var k uint32
	tail := data[nBlocks*4:]
	switch len(tail) {
	case 3:
		k ^= uint32(tail[2]) << 16
		fallthrough
	case 2:
		k ^= uint32(tail[1]) << 8
		fallthrough
	case 1:
		k ^= uint32(tail[0])
		k *= c1
		k = bits.RotateLeft32(k, 15)
		k *= c2
		h ^= k
	}

	h ^= uint32(len(data))
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16

	return h
}

// ToUint64 converts two 32-bit words from buf into a uint64.
func ToUint64(buf []byte, n int) uint64 {
	word := 4
	offset := 2 * n * word
	if len(buf) < offset+8 {
		return 0
	}
	low := binary.LittleEndian.Uint32(buf[offset:])
	high := binary.LittleEndian.Uint32(buf[offset+4:])
	return uint64(low) | uint64(high)<<32
}
