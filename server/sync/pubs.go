package sync

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
)

// Pub exposes the pubs this node knows about,
// and lets other pubs register
type PubsResource struct {
	db *app.DB
}

// NewPubsResource builds a debug
func NewPubsResource(db *app.DB) *PubsResource {
	return &PubsResource{db: db}
}

// GetPubs returns the current user's feed
func (p *PubsResource) GetPubs(w http.ResponseWriter, r *http.Request) {
	pubs, err := p.db.GetPubs()
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(pubs)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}
