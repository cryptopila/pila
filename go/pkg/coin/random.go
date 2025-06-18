package coin

import (
	"crypto/rand"
	"encoding/binary"
	mrand "math/rand"
	"time"
)

func randomUint64() uint64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		// fallback to math/rand seeded with time
		mrand.Seed(time.Now().UnixNano())
		return mrand.Uint64()
	}
	return binary.LittleEndian.Uint64(b[:])
}

// RandomUint8 returns a random uint8 in the range [0,max).
func RandomUint8(max uint8) uint8 {
	if max == 0 {
		return 0
	}
	return uint8(randomUint64() % uint64(max))
}

// RandomUint16 returns a random uint16 in the range [0,max).
func RandomUint16(max uint16) uint16 {
	if max == 0 {
		return 0
	}
	return uint16(randomUint64() % uint64(max))
}

// RandomUint16Range returns a random uint16 between low and high inclusive.
func RandomUint16Range(low, high uint16) uint16 {
	if low >= high {
		return low
	}
	return low + RandomUint16(high-low+1)
}

// RandomUint32Range returns a random uint32 between low and high inclusive.
func RandomUint32Range(low, high uint32) uint32 {
	if low >= high {
		return low
	}
	return low + RandomUint32(high-low+1)
}

// RandomUint32 returns a random uint32 in the range [0,max).
func RandomUint32(max uint32) uint32 {
	if max == 0 {
		return 0
	}
	return uint32(randomUint64() % uint64(max))
}

// RandomUint64 returns a random uint64 in the range [0,max).
func RandomUint64(max uint64) uint64 {
	if max == 0 {
		return 0
	}
	return randomUint64() % max
}

// OpenSSLRANDAdd is a stub to mimic the entropy gathering of the C++ code.
func OpenSSLRANDAdd() {
	// crypto/rand is automatically seeded by the OS; nothing needed
}
