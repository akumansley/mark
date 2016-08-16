package server

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"zenhack.net/go/sandstorm/capnp/hacksession"
	"zenhack.net/go/sandstorm/capnp/sandstormhttpbridge"
	"zombiezen.com/go/capnproto2/rpc"

	"github.com/NYTimes/gziphandler"
	"github.com/PuerkitoBio/goquery"
	"github.com/awans/mark/app"
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

// SandstormHandler reads session info from requests
type SandstormHandler struct {
	i http.Handler
}

func (s SandstormHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HI")
	sessionID := r.Header.Get("X-Sandstorm-Session-Id")
	conn, err := net.Dial("unix", "/tmp/sandstorm-api")
	if err != nil {
		panic(err)
	}
	transport := rpc.StreamTransport(conn)
	ctx := context.Background()
	clientConn := rpc.NewConn(transport)
	defer clientConn.Close()

	bridge := sandstormhttpbridge.SandstormHttpBridge{Client: clientConn.Bootstrap(ctx)}
	fmt.Printf("bridge: %v\n", bridge)
	call := bridge.GetSessionContext(ctx, func(p sandstormhttpbridge.SandstormHttpBridge_getSessionContext_Params) error {
		p.SetId(sessionID)
		return nil
	})
	result, err := call.Struct()
	fmt.Printf("result: %v\n", bridge)
	if err != nil {
		panic(err)
	}
	sc := result.Context()
	hsc := hacksession.HackSessionContext{Client: sc.Client}
	getCall := hsc.HttpGet(ctx, func(p hacksession.HackSessionContext_httpGet_Params) error {
		p.SetUrl("http://www.google.com")
		return nil
	})
	getResult, err := getCall.Struct()
	if err != nil {
		panic(err)
	}
	bytes, err := getResult.Content()

	fmt.Printf("%s\n", bytes)
	// call inner
	s.i.ServeHTTP(w, r)
}

// New returns a new mark server
func New(db *app.DB) http.Handler {
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
	syncRouter.HandleFunc("/announce", a.PutAnnouncement).Methods("POST")

	r.Handle("/bundle.js", http.FileServer(http.Dir("server/data/static/build")))
	r.HandleFunc("/{path:.*}", IndexHandler).Methods("GET")

	gz := gziphandler.GzipHandler(r)
	ss := SandstormHandler{i: gz}
	return ss
}
