package database

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"pila/pkg/coin"
)

// DB wraps a LevelDB instance.
type DB struct {
	db *leveldb.DB
}

// Open opens a LevelDB database located at path.
func Open(path string) (*DB, error) {
	d, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db: d}, nil
}

// Close closes the underlying database.
func (d *DB) Close() error { return d.db.Close() }

// Put stores an arbitrary value under the given key.
func (d *DB) Put(key string, val []byte) error {
	return d.db.Put([]byte(key), val, nil)
}

// Get retrieves the raw value for key.
func (d *DB) Get(key string) ([]byte, error) {
	return d.db.Get([]byte(key), nil)
}

// PutBlock serializes and stores the block using its hash as the key.
func (d *DB) PutBlock(b coin.Block) error {
	if err := b.Validate(); err != nil {
		return err
	}
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return d.Put("block:"+b.Header.Hash(), data)
}

// GetBlock loads the block identified by hash.
func (d *DB) GetBlock(hash string) (coin.Block, error) {
	var out coin.Block
	raw, err := d.Get("block:" + hash)
	if err != nil {
		return out, err
	}
	if err = json.Unmarshal(raw, &out); err != nil {
		return out, err
	}
	if err = out.Validate(); err != nil {
		return out, err
	}
	return out, nil
}

// ListBlocks returns all blocks stored in the database. Any invalid block
// encountered during iteration results in an error.
func (d *DB) ListBlocks() ([]coin.Block, error) {
	var blocks []coin.Block
	iter := d.db.NewIterator(util.BytesPrefix([]byte("block:")), nil)
	for iter.Next() {
		var b coin.Block
		if err := json.Unmarshal(iter.Value(), &b); err != nil {
			iter.Release()
			return nil, err
		}
		if err := b.Validate(); err != nil {
			iter.Release()
			return nil, err
		}
		blocks = append(blocks, b)
	}
	if err := iter.Error(); err != nil {
		iter.Release()
		return nil, err
	}
	iter.Release()
	return blocks, nil
}
