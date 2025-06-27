package crypto

import "testing"

func TestECDHE(t *testing.T) {
	a, err := NewECDHE()
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewECDHE()
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewECDHE()
	if err != nil {
		t.Fatal(err)
	}

	sa1, err := a.Derive(b.Public())
	if err != nil {
		t.Fatalf("derive ab: %v", err)
	}
	sa2, err := b.Derive(a.Public())
	if err != nil {
		t.Fatalf("derive ba: %v", err)
	}
	sc, err := a.Derive(c.Public())
	if err != nil {
		t.Fatalf("derive ac: %v", err)
	}
	if len(sa1) == 0 || len(sa2) == 0 || len(sc) == 0 {
		t.Fatalf("expected secrets")
	}
	if string(sa1) != string(sa2) {
		t.Fatalf("ab secrets mismatch")
	}
	if string(sa1) == string(sc) {
		t.Fatalf("expected different secret for c")
	}
}
