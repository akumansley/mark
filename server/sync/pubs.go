package api

import (
	"net/http"

	"github.com/awans/mark/app"
)

// Pub exposes the pubs this node knows about,
// and lets other pubs register
type PubsResource struct {
	db *app.DB
}

// NewDebug builds a debug
func NewPubsResource(db *app.DB) *Pub {
	return &PubsResource{db: db}
}

// GetDebug returns the current user's feed
func (p *PubsResource) GetPubs(w http.ResponseWriter, r *http.Request) {
	pubs, err := d.db.GetPubs()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}
