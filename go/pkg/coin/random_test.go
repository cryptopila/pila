package coin

import "testing"

func TestRandomRanges(t *testing.T) {
	for i := 0; i < 100; i++ {
		if v := RandomUint8(10); v >= 10 {
			t.Fatalf("uint8 over limit: %d", v)
		}
		if v := RandomUint16Range(5, 10); v < 5 || v > 10 {
			t.Fatalf("uint16 range out: %d", v)
		}
		if v := RandomUint32Range(100, 200); v < 100 || v > 200 {
			t.Fatalf("uint32 range out: %d", v)
		}
	}
}
