package coin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreatePathAndCopyFile(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "a", "b")
	if err := CreatePath(target); err != nil {
		t.Fatalf("create path: %v", err)
	}
	if fi, err := os.Stat(target); err != nil || !fi.IsDir() {
		t.Fatalf("path not created")
	}
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(target, "dst.txt")
	if err := os.WriteFile(src, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("copy file: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil || string(data) != "hello" {
		t.Fatalf("copy failed: %v %s", err, data)
	}
}

func TestPathContents(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f1"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "f2"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	names, err := PathContents(dir)
	if err != nil {
		t.Fatalf("contents: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(names))
	}
}

func TestDataPath(t *testing.T) {
	if DataPath() == "" {
		t.Fatal("empty data path")
	}
}

func TestByteReverse(t *testing.T) {
	if ByteReverse(0x12345678) != 0x78563412 {
		t.Fatalf("byte reverse wrong")
	}
}

func TestDifficultyFromBits(t *testing.T) {
	if DifficultyFromBits(0x1d00ffff) != 1 {
		t.Fatalf("unexpected difficulty")
	}
}

func TestFormatMoney(t *testing.T) {
	if FormatMoney(1234567, false) != "1.234567" {
		t.Fatalf("money format wrong")
	}
	if FormatMoney(Coin, false) != "1.0" {
		t.Fatalf("money format 1 coin")
	}
}
