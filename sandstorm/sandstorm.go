package sandstorm

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"golang.org/x/net/context"
	"zenhack.net/go/sandstorm/capnp/hacksession"
	"zenhack.net/go/sandstorm/capnp/sandstormhttpbridge"
	"zombiezen.com/go/capnproto2/rpc"
)

// Getter implements feed.Getter using the sandstorm-http-bridge
type Getter struct {
	sessionIDs chan string
	m          sync.RWMutex
	sessionID  string
}

// NewGetter returns a new sandstorm.Getter
func NewGetter(sessionIDs chan string) *Getter {
	g := &Getter{sessionIDs: sessionIDs}

	go func() {
		for {
			s := <-sessionIDs
			g.m.Lock()
			g.sessionID = s
			g.m.Unlock()
		}
	}()

	return g
}

// Get implements an HTTP get via sandstorm-http-bridge
func (g *Getter) Get(url string) (*http.Response, error) {
	// block until we have at least one value
	var sessionID string
	for {
		g.m.RLock()
		sessionID = g.sessionID
		g.m.RUnlock()
		if sessionID != "" {
			break
		}
	}

	// TODO how much of this can we re-use across requests..
	conn, err := net.Dial("unix", "/tmp/sandstorm-api")
	if err != nil {
		return nil, err
	}
	transport := rpc.StreamTransport(conn)
	ctx := context.Background()
	clientConn := rpc.NewConn(transport)
	defer clientConn.Close()

	bridge := sandstormhttpbridge.SandstormHttpBridge{Client: clientConn.Bootstrap(ctx)}
	call := bridge.GetSessionContext(ctx, func(p sandstormhttpbridge.SandstormHttpBridge_getSessionContext_Params) error {
		p.SetId(sessionID)
		return nil
	})
	result, err := call.Struct()
	if err != nil {
		return nil, err
	}
	sc := result.Context()
	hsc := hacksession.HackSessionContext{Client: sc.Client}
	getCall := hsc.HttpGet(ctx, func(p hacksession.HackSessionContext_httpGet_Params) error {
		p.SetUrl(url)
		return nil
	})

	// crappy mapping to an http.Response
	ssResponse, err := getCall.Struct()
	if err != nil {
		return nil, err
	}
	content, err := ssResponse.Content() // []byte, err
	if err != nil {
		return nil, err
	}
	mimeType, err := ssResponse.MimeType() // string err
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(content)
	bufCloser := ioutil.NopCloser(buf)

	header := http.Header{}
	header.Add("Content-Type", mimeType)
	res := http.Response{Body: bufCloser, Header: header}
	return &res, nil
}

// Handler reads session info from requests
type Handler struct {
	i http.Handler
	c chan string
}

// NewHandler returns a new sandstorm.Handler and a channel for sesison ids
func NewHandler(inner http.Handler) (http.Handler, chan string) {
	c := make(chan string, 50)
	return Handler{i: inner, c: c}, c
}

func (s Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Sandstorm-Session-Id")
	s.c <- sessionID
	s.i.ServeHTTP(w, r)
}
