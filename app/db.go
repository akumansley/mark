package app

import (
	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"
)

// DB is the application-level DB interface
type DB struct {
	e *entities.DB
}

// NewDB makes a new app db from an entity db
func NewDB(e *entities.DB) *DB {
	return &DB{e: e}
}

// Close closes the underlying db
func (db *DB) Close() {
	db.e.Close()
}

// GetStream returns a user's stream
func (db *DB) GetStream() ([]Bookmark, error) {
	var bookmarks []Bookmark
	db.e.GetAll(&bookmarks)
	return bookmarks, nil
}

// AddBookmark inserts a bookmark into the db
func (db *DB) AddBookmark(b *Bookmark) error {
	id, err := db.e.Add(b)
	b.ID = id
	return err
}

func (db *DB) GetPubs() ([]feed.Pub, error) {
	return db.e.GetPubs()
}

func (db *DB) GetFeeds() ([]feed.SignedFeed, error) {
	return db.e.GetFeeds()
}

func (db *DB) GetFeed(id string) (feed.SignedFeed, error) {
	return db.e.GetFeed(id)
}
