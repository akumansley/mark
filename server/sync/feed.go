package sync

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
	"github.com/gorilla/mux"
)

// FeedResource exposes a feed
type FeedResource struct {
	db *app.DB
}

// NewHeadsResource constructs a HeadsResource
func NewFeedResource(db *app.DB) *FeedResource {
	return &FeedResource{db: db}
}

// GetHeads returns head of each feed
func (f *FeedResource) GetFeed(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	feed, err := f.db.GetFeed(id)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(feed)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}
