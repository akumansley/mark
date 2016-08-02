package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
	p, err := m.db.GetUserProfile()
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(p)
}

// PutProfile updates a user's profile
func (m *Me) PutProfile(w http.ResponseWriter, r *http.Request) {
	var profile app.Profile
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &profile); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	err = m.db.SetProfile(&profile)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		panic(err)
	}
}
