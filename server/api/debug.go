package api

import (
	"crypto/rsa"
	"net/http"

	"github.com/awans/mark"
)

// Debug prints the feed
type Debug struct {
	db  *mark.DB
	key *rsa.PrivateKey
}

// NewDebug builds a debug
func NewDebug(db *mark.DB, key *rsa.PrivateKey) *Debug {
	return &Debug{db: db, key: key}
}

// GetDebug returns the current user's feed
func (d *Debug) GetDebug(w http.ResponseWriter, r *http.Request) {
	feed, err := d.db.FeedForPubKey(&d.key.PublicKey)
	if err != nil {
		panic(err)
	}
	serialized, err := feed.ToBytes(d.key)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(serialized)
}
