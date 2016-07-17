package entities

import (
	"bytes"
	"io"
	"path"

	"github.com/cznic/kv"
)

// Store is a place to keep feed or node data
type Store interface {
	Close() error
	Get([]byte) ([]byte, error)
	Set([]byte, []byte) error
	Delete([]byte) error
	Prefix([]byte) (Iterator, error)
}

// Iterator iterates through keys, returning io.EOF when it's exhausted
type Iterator interface {
	Next() ([]byte, []byte, error)
}

// KvStore is an implementation of Store based on cznic kv
type KvStore struct {
	db *kv.DB
}

var opts = &kv.Options{Compare: nil}

// CreateStore makes a db file for a KvStore
func CreateStore(dirname string) (store Store, err error) {
	filename := path.Join(dirname, "db")
	db, err := kv.Create(filename, opts)
	if err != nil {
		return nil, err
	}

	return KvStore{db: db}, nil
}

// OpenStore opens an existing KvStore
func OpenStore(dirname string) (store Store, err error) {
	filename := path.Join(dirname, "db")

	db, err := kv.Open(filename, opts)
	if err != nil {
		return nil, err
	}

	return KvStore{db: db}, nil
}

// Close closes a KvStore and releases the file
func (kv KvStore) Close() error {
	return kv.db.Close()
}

// Get reads a value from a KvStore
func (kv KvStore) Get(key []byte) ([]byte, error) {
	return kv.db.Get(nil, key)
}

// Set sets a value in a kvstore
func (kv KvStore) Set(key []byte, val []byte) error {
	return kv.db.Set(key, val)
}

// Delete removes a value
func (kv KvStore) Delete(key []byte) error {
	return kv.db.Delete(key)
}

type kvIterator struct {
	e      *kv.Enumerator
	prefix []byte
}

// Prefix implements Store
func (kv KvStore) Prefix(key []byte) (Iterator, error) {
	e, _, err := kv.db.Seek(key)
	if err != nil {
		return nil, err
	}
	return kvIterator{e: e, prefix: key}, nil
}

// Next implements Iterator
func (i kvIterator) Next() ([]byte, []byte, error) {
	k, v, err := i.e.Next()
	if err != nil {
		return nil, nil, err
	}
	if bytes.HasPrefix(k, i.prefix) {
		return k, v, nil
	}
	return nil, nil, io.EOF
}
