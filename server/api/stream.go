package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/awans/mark/app"
)

// Stream (not a Feed) represents bookmarks across users
type Stream struct {
	db *app.DB
}

// NewStream builds a stream resource
func NewStream(db *app.DB) *Stream {
	return &Stream{db: db}
}

type streamBookmark struct {
	app.Bookmark
	Profile *app.Profile `json:"profile"`
}

// GetStream returns the current user's stream
func (s *Stream) GetStream(w http.ResponseWriter, r *http.Request) {
	countS := r.URL.Query()["count"][0]
	count, err := strconv.Atoi(countS)
	if err != nil {
		panic(err)
	}
	offsetS := r.URL.Query()["offset"][0]
	offset, err := strconv.Atoi(offsetS)
	if err != nil {
		panic(err)
	}

	bookmarks, err := s.db.GetStream(count, offset)
	if err != nil {
		panic(err)
	}
	if bookmarks == nil {
		bookmarks = make([]app.Bookmark, 0)
	}

	sbs := make([]streamBookmark, 0)
	for _, b := range bookmarks {
		p, err := s.db.GetProfile(b.FeedID)
		if err != nil {
			panic(err)
		}
		sbs = append(sbs, streamBookmark{Bookmark: b, Profile: p})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(sbs)
}
