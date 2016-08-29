package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/awans/mark/app"
	"github.com/awans/mark/feed"
)

// Self represents the current server
type Self struct {
	db *app.DB
}

// NewSelf builds a self resource
func NewSelf(db *app.DB) *Self {
	return &Self{db: db}
}

// GetSelf returns the current self pub
func (s *Self) GetSelf(w http.ResponseWriter, r *http.Request) {
	self, err := s.db.GetSelf()
	if err != nil {
		panic(err)
	}
	if self == nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(self)
}

// PutSelf lets a user update the self url
func (s *Self) PutSelf(w http.ResponseWriter, r *http.Request) {
	var self feed.Pub
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &self); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	err = s.db.PutSelf(&self)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(self); err != nil {
		panic(err)
	}
}
