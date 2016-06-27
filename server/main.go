package server

import "net/http"
import "github.com/gorilla/mux"

// AppHandler handles all requests that want to return the client SPA
func AppHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/index.html")
}

// APIHandler handles all requests that want an api response
func APIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  res := vars["resource"]
	w.Write([]byte(res))
}

// New returns a new mark server
func New() http.Handler {
	r := mux.NewRouter()

	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/{resource}", APIHandler)

	r.Handle("/{path:.*}", http.FileServer(http.Dir("server/data/static/build")))
	return r
}
