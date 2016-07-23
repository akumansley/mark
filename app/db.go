package app

import (
	"crypto/rsa"

	"github.com/awans/mark/entities"
)

// DB is the application-level DB interface
type DB struct {
	key *rsa.PrivateKey // still maybe hide this in a Session
	e   *entities.DB
}

// NewDB makes a new app db from an entity db
func NewDB(e *entities.DB, key *rsa.PrivateKey) *DB {
	return &DB{e: e, key: key}
}

// Close closes the underlying db
func (db *DB) Close() {
	db.e.Close()
}

// GetFeed returns a user's feed
func (db *DB) GetFeed() ([]Bookmark, error) {
	var bookmarks []Bookmark
	db.e.GetAll(&bookmarks)
	return bookmarks, nil
}

// AddBookmark inserts a bookmark into the db
func (db *DB) AddBookmark(b *Bookmark) {
	id, _ := db.e.Add(b)
	b.ID = id
}

// DebugFeed returns the user's feed
func (db *DB) DebugFeed() ([]byte, error) {
	feed, err := db.e.UserFeed()
	if err != nil {
		return nil, err
	}
	return feed.ToBytes(db.key)
}
