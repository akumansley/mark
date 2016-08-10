package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/awans/mark/app"
	"github.com/gorilla/mux"
)

// Bookmark is a resource
type Bookmark struct {
	db *app.DB
}

// NewBookmark builds a bookmark resource
func NewBookmark(db *app.DB) *Bookmark {
	return &Bookmark{db: db}
}

// AddBookmark creates a new bookmark
func (b *Bookmark) AddBookmark(w http.ResponseWriter, r *http.Request) {
	var bookmark app.Bookmark
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &bookmark); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	err = b.db.AddBookmark(&bookmark)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bookmark); err != nil {
		panic(err)
	}
}

// RemoveBookmark removes an existing bookmark
func (b *Bookmark) RemoveBookmark(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := b.db.RemoveBookmark(id)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}
