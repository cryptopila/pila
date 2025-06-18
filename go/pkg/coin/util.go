package coin

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"syscall"

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

// Abs64 returns the absolute value of v as an int64.
func Abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

// MoneyRange returns true if value is within the allowed monetary range.
func MoneyRange(value int64) bool { return value >= 0 && value <= MaxMoneySupply }

// FormatMoney converts the amount to a human readable decimal string. If plus
// is true a leading '+' will be included for positive values.
func FormatMoney(n int64, plus bool) string {
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	} else if plus && n > 0 {
		sign = "+"
	}

	quotient := n / Coin
	remainder := n % Coin
	s := fmt.Sprintf("%s%d.%06d", sign, quotient, remainder)
	s = strings.TrimRight(s, "0")
	if strings.HasSuffix(s, ".") {
		s += "0"
	}
	return s
}

// FormatVersion converts the integer protocol version to dotted form.
func FormatVersion(version int32) string {
	if version%100 == 0 {
		return fmt.Sprintf("%d.%d.%d",
			version/1000000,
			(version/10000)%100,
			(version/100)%100)
	}
	return fmt.Sprintf("%d.%d.%d.%d",
		version/1000000,
		(version/10000)%100,
		(version/100)%100,
		version%100)
}

// FormatSubVersion returns a BIP0014 compliant subversion string.
func FormatSubVersion(name string, clientVersion int32, comments []string) string {
	out := "/" + name + ":" + FormatVersion(clientVersion)
	if len(comments) > 0 {
		out += "(" + strings.Join(comments, "; ") + ")"
	}
	return out + "/"
}

// HexString converts bytes into a hex string. If spaces is true the bytes are
// separated by spaces.
func HexString(data []byte, spaces bool) string {
	if !spaces {
		return hex.EncodeToString(data)
	}
	parts := make([]string, len(data))
	for i, b := range data {
		parts[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(parts, " ")
}

// HexStringFromBits encodes a uint32 as big-endian hex.
func HexStringFromBits(bits uint32) string {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], bits)
	return hex.EncodeToString(buf[:])
}

// IsHex returns true if the string is non-empty and consists of pairs of
// hexadecimal characters.
func IsHex(val string) bool {
	if len(val) == 0 || len(val)%2 != 0 {
		return false
	}
	for i := 0; i < len(val); i++ {
		c := val[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// FromHex converts a hex string into a byte slice, ignoring spaces.
func FromHex(val string) []byte {
	val = strings.ReplaceAll(val, " ", "")
	if len(val)%2 != 0 {
		return nil
	}
	out := make([]byte, len(val)/2)
	_, _ = hex.Decode(out, []byte(val))
	return out
}

// GetVarIntSize returns the size of a variable integer encoding of n.
func GetVarIntSize(n uint64) uint32 {
	switch {
	case n < 253:
		return 1
	case n <= math.MaxUint16:
		return 1 + 2
	case n <= math.MaxUint32:
		return 1 + 4
	default:
		return 1 + 8
	}
}

// ByteReverse swaps the byte order of a 32-bit word.
func ByteReverse(val uint32) uint32 {
	ret := ((val & 0xFF00FF00) >> 8) | ((val & 0x00FF00FF) << 8)
	return (ret << 16) | (ret >> 16)
}

// DifficultyFromBits converts compact representation to a difficulty value.
func DifficultyFromBits(bits uint32) float64 {
	shift := (bits >> 24) & 0xff
	diff := float64(0x0000ffff) / float64(bits&0x00ffffff)
	for shift < 29 {
		diff *= 256.0
		shift++
	}
	for shift > 29 {
		diff /= 256.0
		shift--
	}
	return diff
}

// DiskInfo contains filesystem statistics returned by DiskInfo.
type DiskInfo struct {
	Capacity  uint64
	Free      uint64
	Available uint64
}

// DiskInfo obtains disk usage statistics for the given path.
func DiskInfoPath(path string) DiskInfo {
	var fs syscall.Statfs_t
	if err := syscall.Statfs(path, &fs); err == nil {
		cap := fs.Blocks * uint64(fs.Bsize)
		free := fs.Bfree * uint64(fs.Bsize)
		avail := fs.Bavail * uint64(fs.Bsize)
		return DiskInfo{cap, free, avail}
	}
	return DiskInfo{}
}
