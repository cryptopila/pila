package coin

import (
	"crypto/rand"
	"encoding/binary"
	"math/bits"

	"github.com/jzelinskie/whirlpool"
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

// WhirlpoolX computes a 32-byte digest derived from the whirlpool hash.
// It mirrors the C++ implementation which XORs the first 32 bytes of the
// digest with a 16-byte shifted portion of the same digest.
func WhirlpoolX(data []byte) [32]byte {
	h := whirlpool.New()
	h.Write(data)
	full := h.Sum(nil)
	var out [32]byte
	for i := 0; i < 32; i++ {
		out[i] = full[i] ^ full[i+16]
	}
	return out
}

// Blake256EightRound computes the 8-round BLAKE-256 hash used by the
// original C++ code.  This implementation is self-contained and processes
// the input in 512-bit blocks following the standard padding rules.
func Blake256EightRound(data []byte) [32]byte {
	var h = [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
		0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}
	var s [4]uint32
	var counter uint64
	var block [64]byte

	// Process full 64-byte blocks.
	for len(data) >= 64 {
		copy(block[:], data[:64])
		counter += 512
		compress8(&h, &s, block[:], counter)
		data = data[64:]
	}

	// Finalize with padding per the BLAKE specification.
	msgBits := counter + uint64(len(data))<<3
	if len(data) < 55 {
		copy(block[:], data)
		block[len(data)] = 0x80
		for i := len(data) + 1; i < 55; i++ {
			block[i] = 0
		}
		block[55] = 0x01
		binary.BigEndian.PutUint64(block[56:], msgBits)
		compress8(&h, &s, block[:], msgBits)
	} else {
		copy(block[:], data)
		block[len(data)] = 0x80
		for i := len(data) + 1; i < 64; i++ {
			block[i] = 0
		}
		counter += uint64(64-len(data)) << 3
		compress8(&h, &s, block[:], counter)

		for i := 0; i < 55; i++ {
			block[i] = 0
		}
		block[55] = 0x01
		binary.BigEndian.PutUint64(block[56:], msgBits)
		compress8(&h, &s, block[:], 0)
	}

	var out [32]byte
	binary.BigEndian.PutUint32(out[28:], h[7])
	binary.BigEndian.PutUint32(out[24:], h[6])
	binary.BigEndian.PutUint32(out[20:], h[5])
	binary.BigEndian.PutUint32(out[16:], h[4])
	binary.BigEndian.PutUint32(out[12:], h[3])
	binary.BigEndian.PutUint32(out[8:], h[2])
	binary.BigEndian.PutUint32(out[4:], h[1])
	binary.BigEndian.PutUint32(out[0:], h[0])
	return out
}

// g is the quarter round function for BLAKE-256.
func g(a, b, c, d, mx, my, cx, cy uint32) (uint32, uint32, uint32, uint32) {
	a += b + (mx ^ cy)
	d = bits.RotateLeft32(d^a, -16)
	c += d
	b = bits.RotateLeft32(b^c, -12)
	a += b + (my ^ cx)
	d = bits.RotateLeft32(d^a, -8)
	c += d
	b = bits.RotateLeft32(b^c, -7)
	return a, b, c, d
}

// compress8 performs the 8-round BLAKE-256 compression on a single block.
func compress8(h *[8]uint32, s *[4]uint32, block []byte, counter uint64) {
	const (
		c0 uint32 = 0x243f6a88
		c1 uint32 = 0x85a308d3
		c2 uint32 = 0x13198a2e
		c3 uint32 = 0x03707344
		c4 uint32 = 0xa4093822
		c5 uint32 = 0x299f31d0
		c6 uint32 = 0x082efa98
		c7 uint32 = 0xec4e6c89
		c8 uint32 = 0x452821e6
		c9 uint32 = 0x38d01377
		ca uint32 = 0xbe5466cf
		cb uint32 = 0x34e90c6c
		cc uint32 = 0xc0ac29b7
		cd uint32 = 0xc97c50dd
		ce uint32 = 0x3f84d5b5
		cf uint32 = 0xb5470917
	)

	m0 := binary.BigEndian.Uint32(block[0:4])
	m1 := binary.BigEndian.Uint32(block[4:8])
	m2 := binary.BigEndian.Uint32(block[8:12])
	m3 := binary.BigEndian.Uint32(block[12:16])
	m4 := binary.BigEndian.Uint32(block[16:20])
	m5 := binary.BigEndian.Uint32(block[20:24])
	m6 := binary.BigEndian.Uint32(block[24:28])
	m7 := binary.BigEndian.Uint32(block[28:32])
	m8 := binary.BigEndian.Uint32(block[32:36])
	m9 := binary.BigEndian.Uint32(block[36:40])
	m10 := binary.BigEndian.Uint32(block[40:44])
	m11 := binary.BigEndian.Uint32(block[44:48])
	m12 := binary.BigEndian.Uint32(block[48:52])
	m13 := binary.BigEndian.Uint32(block[52:56])
	m14 := binary.BigEndian.Uint32(block[56:60])
	m15 := binary.BigEndian.Uint32(block[60:64])

	t0 := uint32(counter)
	t1 := uint32(counter >> 32)

	v0, v4, v8, vc := g(h[0], h[4], s[0]^c0, t0^c4, m0, m1, c1, c0)
	v1, v5, v9, vd := g(h[1], h[5], s[1]^c1, t0^c5, m2, m3, c3, c2)
	v2, v6, va, ve := g(h[2], h[6], s[2]^c2, t1^c6, m4, m5, c5, c4)
	v3, v7, vb, vf := g(h[3], h[7], s[3]^c3, t1^c7, m6, m7, c7, c6)
	v0, v5, va, vf = g(v0, v5, va, vf, m8, m9, c9, c8)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m10, m11, cb, ca)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m12, m13, cd, cc)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m14, m15, cf, ce)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m14, m10, ca, ce)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m4, m8, c8, c4)
	v2, v6, va, ve = g(v2, v6, va, ve, m9, m15, cf, c9)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m13, m6, c6, cd)
	v0, v5, va, vf = g(v0, v5, va, vf, m1, m12, cc, c1)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m0, m2, c2, c0)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m11, m7, c7, cb)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m5, m3, c3, c5)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m11, m8, c8, cb)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m12, m0, c0, cc)
	v2, v6, va, ve = g(v2, v6, va, ve, m5, m2, c2, c5)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m15, m13, cd, cf)
	v0, v5, va, vf = g(v0, v5, va, vf, m10, m14, ce, ca)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m3, m6, c6, c3)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m7, m1, c1, c7)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m9, m4, c4, c9)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m7, m9, c9, c7)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m3, m1, c1, c3)
	v2, v6, va, ve = g(v2, v6, va, ve, m13, m12, cc, cd)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m11, m14, ce, cb)
	v0, v5, va, vf = g(v0, v5, va, vf, m2, m6, c6, c2)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m5, m10, ca, c5)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m4, m0, c0, c4)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m15, m8, c8, cf)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m9, m0, c0, c9)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m5, m7, c7, c5)
	v2, v6, va, ve = g(v2, v6, va, ve, m2, m4, c4, c2)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m10, m15, cf, ca)
	v0, v5, va, vf = g(v0, v5, va, vf, m14, m1, c1, ce)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m11, m12, cc, cb)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m6, m8, c8, c6)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m3, m13, cd, c3)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m2, m12, cc, c2)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m6, m10, ca, c6)
	v2, v6, va, ve = g(v2, v6, va, ve, m0, m11, cb, c0)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m8, m3, c3, c8)
	v0, v5, va, vf = g(v0, v5, va, vf, m4, m13, cd, c4)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m7, m5, c5, c7)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m15, m14, ce, cf)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m1, m9, c9, c1)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m12, m5, c5, cc)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m1, m15, cf, c1)
	v2, v6, va, ve = g(v2, v6, va, ve, m14, m13, cd, ce)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m4, m10, ca, c4)
	v0, v5, va, vf = g(v0, v5, va, vf, m0, m7, c7, c0)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m6, m3, c3, c6)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m9, m2, c2, c9)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m8, m11, cb, c8)

	v0, v4, v8, vc = g(v0, v4, v8, vc, m13, m11, cb, cd)
	v1, v5, v9, vd = g(v1, v5, v9, vd, m7, m14, ce, c7)
	v2, v6, va, ve = g(v2, v6, va, ve, m12, m1, c1, cc)
	v3, v7, vb, vf = g(v3, v7, vb, vf, m3, m9, c9, c3)
	v0, v5, va, vf = g(v0, v5, va, vf, m5, m0, c0, c5)
	v1, v6, vb, vc = g(v1, v6, vb, vc, m15, m4, c4, cf)
	v2, v7, v8, vd = g(v2, v7, v8, vd, m8, m6, c6, c8)
	v3, v4, v9, ve = g(v3, v4, v9, ve, m2, m10, ca, c2)

	h[0] ^= s[0] ^ v0 ^ v8
	h[1] ^= s[1] ^ v1 ^ v9
	h[2] ^= s[2] ^ v2 ^ va
	h[3] ^= s[3] ^ v3 ^ vb
	h[4] ^= s[0] ^ v4 ^ vc
	h[5] ^= s[1] ^ v5 ^ vd
	h[6] ^= s[2] ^ v6 ^ ve
	h[7] ^= s[3] ^ v7 ^ vf
}
