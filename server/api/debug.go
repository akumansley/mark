package api

import (
	"net/http"

	"github.com/awans/mark/app"
)

// Debug prints the feed
type Debug struct {
	db  *app.DB
}

// NewDebug builds a debug
func NewDebug(db *app.DB) *Debug {
	return &Debug{db: db}
}

// GetDebug returns the current user's feed
func (d *Debug) GetDebug(w http.ResponseWriter, r *http.Request) {
	bytes, err := d.db.DebugFeed()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}
