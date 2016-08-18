package sync

import (
	"net/http"

	"github.com/awans/mark/app"
	"github.com/awans/mark/feed"
)

// AnnounceResource accepts announcements
type AnnounceResource struct {
	db *app.DB
}

// NewAnnounceResource constructs a AnnounceResource
func NewAnnounceResource(db *app.DB) *AnnounceResource {
	return &AnnounceResource{db: db}
}

// GetAnnouncement accepts updates to pubs
func (a *AnnounceResource) GetAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcedURL := r.URL.Query().Get("url")

	go func() {
		p := feed.Pub{URL: announcedURL}
		a.db.PutPub(&p)
		feeds, err := a.db.GetFeeds()
		newPubs, newFeeds, err := feed.Sync([]feed.Pub{p}, feeds)
		if err != nil {
			panic(err)
		}
		for _, f := range newFeeds {
			a.db.PutFeed(f)
		}
		for _, p := range newPubs {
			a.db.PutPub(&p)
		}
	}()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
}
