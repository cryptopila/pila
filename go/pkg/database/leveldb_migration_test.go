package database

import (
	"path/filepath"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func openDB(t *testing.T, path string, opts *opt.Options) *leveldb.DB {
	t.Helper()
	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestLevelDBMigration(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testdb")

	t.Run("open missing", func(t *testing.T) {
		if _, err := leveldb.OpenFile(path, &opt.Options{ErrorIfMissing: true}); err == nil {
			t.Fatalf("expected missing db error")
		}
	})

	t.Run("basic ops", func(t *testing.T) {
		db := openDB(t, path, nil)

		if err := db.Put([]byte("foo"), []byte("hello"), nil); err != nil {
			t.Fatalf("put: %v", err)
		}
		val, err := db.Get([]byte("foo"), nil)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if string(val) != "hello" {
			t.Fatalf("expected hello, got %s", val)
		}

		batch := new(leveldb.Batch)
		batch.Put([]byte("bar"), []byte("b"))
		batch.Put([]byte("box"), []byte("c"))
		batch.Delete([]byte("bar"))
		if err := db.Write(batch, nil); err != nil {
			t.Fatalf("write batch: %v", err)
		}

		if _, err := db.Get([]byte("bar"), nil); err != leveldb.ErrNotFound {
			t.Fatalf("expected bar to be deleted")
		}
		val, err = db.Get([]byte("box"), nil)
		if err != nil || string(val) != "c" {
			t.Fatalf("unexpected box value %q, err %v", val, err)
		}

		iter := db.NewIterator(nil, nil)
		if !iter.First() || string(iter.Key()) != "box" {
			t.Fatalf("expected first key box")
		}
		if !iter.Next() || string(iter.Key()) != "foo" {
			t.Fatalf("expected second key foo")
		}
		iter.Release()
		if err := iter.Error(); err != nil {
			t.Fatalf("iter err: %v", err)
		}

		snap, err := db.GetSnapshot()
		if err != nil {
			t.Fatalf("snapshot: %v", err)
		}
		t.Cleanup(func() { snap.Release() })
		if err := db.Delete([]byte("foo"), nil); err != nil {
			t.Fatalf("del: %v", err)
		}
		val, err = snap.Get([]byte("foo"), nil)
		if err != nil || string(val) != "hello" {
			t.Fatalf("snapshot get: %q err %v", val, err)
		}

		if err := db.CompactRange(util.Range{}); err != nil {
			t.Fatalf("compact: %v", err)
		}
		if val, err = db.Get([]byte("box"), nil); err != nil || string(val) != "c" {
			t.Fatalf("post compact: %q err %v", val, err)
		}

		if _, err := db.GetProperty("nosuchprop"); err == nil {
			t.Fatalf("expected invalid property error")
		}
		if prop, err := db.GetProperty("leveldb.stats"); err != nil || prop == "" {
			t.Fatalf("expected stats property, got %q err %v", prop, err)
		}
	})

	t.Run("recover", func(t *testing.T) {
		db, err := leveldb.RecoverFile(path, nil)
		if err != nil {
			t.Fatalf("recover: %v", err)
		}
		t.Cleanup(func() { db.Close() })
		if _, err := db.Get([]byte("foo"), nil); err != leveldb.ErrNotFound {
			t.Fatalf("expected foo to be deleted after recover")
		}
		if val, err := db.Get([]byte("box"), nil); err != nil || string(val) != "c" {
			t.Fatalf("recovered box: %q err %v", val, err)
		}
	})

	t.Run("bloom reopen", func(t *testing.T) {
		db := openDB(t, path, &opt.Options{Filter: filter.NewBloomFilter(10)})

		if err := db.Put([]byte("foo"), []byte("foovalue"), nil); err != nil {
			t.Fatalf("put foo: %v", err)
		}
		if err := db.Put([]byte("bar"), []byte("barvalue"), nil); err != nil {
			t.Fatalf("put bar: %v", err)
		}
		if err := db.CompactRange(util.Range{}); err != nil {
			t.Fatalf("compact with filter: %v", err)
		}
		if val, err := db.Get([]byte("foo"), nil); err != nil || string(val) != "foovalue" {
			t.Fatalf("foo after filter reopen: %q err %v", val, err)
		}
		if val, err := db.Get([]byte("bar"), nil); err != nil || string(val) != "barvalue" {
			t.Fatalf("bar after filter reopen: %q err %v", val, err)
		}
	})
}
