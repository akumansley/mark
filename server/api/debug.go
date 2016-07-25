package api

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
)

// Debug prints the feed
type Debug struct {
	db *app.DB
}

// NewDebug builds a debug
func NewDebug(db *app.DB) *Debug {
	return &Debug{db: db}
}

// GetDebug returns all feeds
func (d *Debug) GetDebug(w http.ResponseWriter, r *http.Request) {
	sfs, err := d.db.GetFeeds()
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(sfs)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}
