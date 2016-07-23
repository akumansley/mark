package api

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
)

// Stream (not a Feed) represents bookmarks across users
type Stream struct {
	db *app.DB
}

// NewFeed builds a stream resource
func NewStream(db *app.DB) *Stream {
	return &Stream{db: db}
}

// GetFeed returns the current user's stream
func (s *Stream) GetStream(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := s.db.GetFeed()
	if err != nil {
		panic(err)
	}
	if bookmarks == nil {
		bookmarks = make([]app.Bookmark, 0)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(bookmarks)
}
