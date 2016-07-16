package api

import (
	"crypto/rsa"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/awans/mark"
)

// Bookmark is a resource
type Bookmark struct {
	db  *mark.DB
	key *rsa.PrivateKey
}

// NewBookmark builds a bookmark resource
func NewBookmark(db *mark.DB, key *rsa.PrivateKey) *Bookmark {
	return &Bookmark{db: db, key: key}
}

// AddBookmark creates a new bookmark
func (b *Bookmark) AddBookmark(w http.ResponseWriter, r *http.Request) {
	var bookmark mark.Bookmark
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
	var bookmarks []mark.Bookmark
	b.db.GetAll(&b.key.PublicKey, &bookmarks)
	entity := mark.Entity{ID: strconv.Itoa(len(bookmarks)), Body: &bookmark}
	err = b.db.Add(b.key, entity)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(bookmark); err != nil {
		panic(err)
	}
}
