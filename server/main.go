package server

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/PuerkitoBio/goquery"
	"github.com/awans/mark/app"
	"github.com/awans/mark/sandstorm"
	"github.com/awans/mark/server/api"
	"github.com/awans/mark/server/sync"
	"github.com/gorilla/mux"
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
func New(db *app.DB) (http.Handler, *sandstorm.SessionBus) {
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()
	s := api.NewStream(db)
	apiRouter.HandleFunc("/stream", s.GetStream).Methods("GET")
	b := api.NewBookmark(db)
	apiRouter.HandleFunc("/bookmark", b.AddBookmark).Methods("POST")
	apiRouter.HandleFunc("/bookmark/{id}", b.RemoveBookmark).Methods("DELETE")
	d := api.NewDebug(db)
	apiRouter.HandleFunc("/debug", d.GetDebug).Methods("GET")
	me := api.NewMe(db)
	apiRouter.HandleFunc("/me", me.GetProfile).Methods("GET")
	apiRouter.HandleFunc("/me", me.PutProfile).Methods("PUT")
	self := api.NewSelf(db)
	apiRouter.HandleFunc("/self", self.GetSelf).Methods("GET")
	apiRouter.HandleFunc("/self", self.PutSelf).Methods("PUT")

	viewsRouter := r.PathPrefix("/views").Subrouter()
	viewsRouter.HandleFunc("/title", TitleHandler).Methods("GET")

	syncRouter := r.PathPrefix("/sync").Subrouter()
	p := sync.NewPubsResource(db)
	syncRouter.HandleFunc("/pubs", p.GetPubs).Methods("GET")
	h := sync.NewHeadsResource(db)
	syncRouter.HandleFunc("/heads", h.GetHeads).Methods("GET")
	f := sync.NewFeedResource(db)
	syncRouter.HandleFunc("/feed/{id}", f.GetFeed).Methods("GET")

	a := sync.NewAnnounceResource(db)
	// Shouldn't really be a get, but is due to limitations of sandstorm
	syncRouter.HandleFunc("/announce", a.GetAnnouncement).Methods("GET")

	r.Handle("/bundle.js", http.FileServer(http.Dir("server/data/static/build")))
	r.HandleFunc("/{path:.*}", IndexHandler).Methods("GET")

	gz := gziphandler.GzipHandler(r)

	// the sandstorm handler intercepts the sandstorm session ID and passes it to the Getter
	// So background requests are made with the sessionID that "last touched" the app
	ss, bus := sandstorm.NewHandler(gz)
	return ss, bus
}
