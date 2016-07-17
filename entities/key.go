package entities

import "bytes"

// Key is a db key
type Key struct {
	path [][]byte
}

// Separator is the path separator for store keys
var Separator = []byte("/")

// NewKey returns a key from a path
func NewKey(components ...string) *Key {
	var path [][]byte
	for _, c := range components {
		path = append(path, []byte(c))
	}
	return &Key{path: path}
}

// ToBytes renders the key to bytes
func (k *Key) ToBytes() []byte {
	return bytes.Join(k.path, Separator)
}
