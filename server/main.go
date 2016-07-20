package server

import (
	"net/http"

	"github.com/awans/mark/app"
	"github.com/awans/mark/server/api"
	"github.com/gorilla/mux"
	"github.com/PuerkitoBio/goquery"
	"github.com/kennygrant/sanitize"

)

// TitleHandler returns the page title of a url
func TitleHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()["url"][0]

	doc, err := goquery.NewDocument(url)
	// TODO
	if err != nil {
		panic(err)
	}
	title := doc.Find("title").Text()
	sanitized := sanitize.HTML(title)
	w.Write([]byte(sanitized))
}

// IndexHandler serves index.html (ie the SPA)
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "server/data/static/build/index.html")
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

	viewsRouter := r.PathPrefix("/views").Subrouter()
	viewsRouter.HandleFunc("/title", TitleHandler).Methods("GET")

	r.Handle("/bundle.js", http.FileServer(http.Dir("server/data/static/build")))
	r.HandleFunc("/{path:.*}", IndexHandler).Methods("GET")
	return r
}
