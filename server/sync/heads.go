package sync

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
	"github.com/awans/mark/feed"
)

// Heads exposes the head of each feed
type HeadsResource struct {
	db *app.DB
}

// NewHeadsResource constructs a HeadsResource
func NewHeadsResource(db *app.DB) *HeadsResource {
	return &HeadsResource{db: db}
}

// GetHeads returns head of each feed
func (h *HeadsResource) GetHeads(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.db.GetFeeds()
	if err != nil {
		panic(err)
	}

	var heads []feed.Head
	for _, f := range feeds {
		fp, err := f.Fingerprint()
		if err != nil {
			panic(err)
		}
		heads = append(heads, feed.Head{ID: string(fp), Len: len(f)})
	}

	bytes, err := json.Marshal(heads)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}
