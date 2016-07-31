package api

import (
	"encoding/json"
	"net/http"

	"github.com/awans/mark/app"
)

// Me represents the signed in user
type Me struct {
	db *app.DB
}

// NewMe builds a profile resource
func NewMe(db *app.DB) *Me {
	return &Me{db: db}
}

// GetProfile returns the current user's profile
func (m *Me) GetProfile(w http.ResponseWriter, r *http.Request) {
	p, err := m.db.GetProfile()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(p)
}
