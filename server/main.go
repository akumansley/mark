package server

import (
	"net/http"

	"github.com/awans/mark/app"
	"github.com/awans/mark/server/api"
	"github.com/gorilla/mux"
)

// AppHandler handles all requests that want to return the client SPA
func AppHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/index.html")
}

// New returns a new mark server
func New(db *app.DB) http.Handler {
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()

	f := api.NewFeed(db)
	apiRouter.HandleFunc("/feed", f.GetFeed).Methods("GET")

	b := api.NewBookmark(db)
	apiRouter.HandleFunc("/bookmark", b.AddBookmark).Methods("POST")

	d := api.NewDebug(db)
	apiRouter.HandleFunc("/debug", d.GetDebug).Methods("GET")

	r.Handle("/{path:.*}", http.FileServer(http.Dir("server/data/static/build")))
	return r
}
