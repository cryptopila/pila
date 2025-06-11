package pkg

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
)

type DB struct {
	inner *leveldb.DB
}

func NewDB(path string) (*DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{inner: db}, nil
}

func (d *DB) Close() error {
	if d.inner == nil {
		return nil
	}
	return d.inner.Close()
}

func (d *DB) PutBlock(b *Block) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return d.inner.Put(b.Hash, data, nil)
}

func (d *DB) GetBlock(hash []byte) (*Block, error) {
	data, err := d.inner.Get(hash, nil)
	if err != nil {
		return nil, err
	}
	var b Block
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
