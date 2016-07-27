package sync

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

// PutAnnouncement accepts updates to pubs and feeds
func (a *AnnounceResource) PutAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announce feed.Announcement

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &announce); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	a.db.PutPub(&announce.Pub)
	feeds, err := a.db.GetFeeds()
	newPubs, newFeeds, err := feed.Sync([]feed.Pub{announce.Pub}, feeds)
	if err != nil {
		panic(err)
	}
	for _, f := range newFeeds {
		a.db.PutFeed(f)
	}
	for _, p := range newPubs {
		a.db.PutPub(&p)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusAccepted)
}
