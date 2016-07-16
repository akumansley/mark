package api

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/awans/mark"
)

// Feed represents bookmarks across users
type Feed struct {
	db  *mark.DB
	key *rsa.PublicKey
}

// NewFeed builds a feed resource
func NewFeed(db *mark.DB, key *rsa.PublicKey) *Feed {
	return &Feed{db: db, key: key}
}

// GetFeed returns the current user's feed
func (f *Feed) GetFeed(w http.ResponseWriter, r *http.Request) {
	var bookmarks []mark.Bookmark
	f.db.GetAll(f.key, &bookmarks)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(bookmarks)
}
