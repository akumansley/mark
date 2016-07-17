package api

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
)

// Feed represents bookmarks across users
type Feed struct {
	db *app.DB
	u  *app.User
}

// NewFeed builds a feed resource
func NewFeed(db *app.DB) *Feed {
	return &Feed{db: db}
}

// GetFeed returns the current user's feed
func (f *Feed) GetFeed(w http.ResponseWriter, r *http.Request) {
	bookmarks, err := f.db.GetFeed()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(bookmarks)
}
