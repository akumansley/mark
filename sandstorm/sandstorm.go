package sandstorm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	"github.com/awans/mark/entities"
	"github.com/awans/mark/feed"

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
func NewGetter(bus *SessionBus) *Getter {
	c := bus.Sub()
	g := &Getter{sessionIDs: c}

	go func() {
		for {
			s := <-c
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

// SessionBus gives us one to many broadcasting of ss sessionIDs
type SessionBus struct {
	listeners []chan string
}

// Pub writes a new sessionID to all channels
func (s *SessionBus) Pub(sessionID string) {
	for _, c := range s.listeners {
		c <- sessionID
	}
}

// Sub returns a new channel with sessionIDs
func (s *SessionBus) Sub() chan string {
	c := make(chan string, 50)
	s.listeners = append(s.listeners, c)
	return c
}

// Handler reads session info from requests
type Handler struct {
	i http.Handler
	b *SessionBus
}

// NewHandler returns a new sandstorm.Handler and a channel for sesison ids
func NewHandler(inner http.Handler) (http.Handler, *SessionBus) {
	b := SessionBus{}
	return Handler{i: inner, b: &b}, &b
}

func (s Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionID := r.Header.Get("X-Sandstorm-Session-Id")
	s.b.Pub(sessionID)
	s.i.ServeHTTP(w, r)
}

// GetEndpointURL will, if it hasn't already been done
// ask sandstorm for an endpoint URL and save it
// in the DB
func GetEndpointURL(db *entities.DB, bus *SessionBus) {
	self, err := db.GetSelf()
	if self != nil {
		fmt.Printf("Got Pub URL %s\n", self.URL)
		return
	}
	if err != nil || self == nil {
		// we don't know where we are on the internet, so let's do it!
		c := bus.Sub()
		// block on a session
		fmt.Println("Blocking on sessionID")
		sessionID := <-c
		fmt.Printf("sessionID %s\n", sessionID)
		conn, err := net.Dial("unix", "/tmp/sandstorm-api")
		if err != nil {
			fmt.Printf("Error getting endpoint URL: %s\n", err)
			return
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
			fmt.Printf("Error getting endpoint URL: %s\n", err)
			return
		}
		sc := result.Context()
		hsc := hacksession.HackSessionContext{Client: sc.Client}
		apiTokenCall := hsc.GenerateApiToken(ctx, func(p hacksession.HackSessionContext_generateApiToken_Params) error {
			// 	generateApiToken @3 (petname :Text, userInfo :Grain.UserInfo, expires :UInt64 = 0)
			//  -> (token :Text, endpointUrl :Text, tokenId :Text);
			p.SetPetname("Mark API") // what's this even do
			return nil
		})
		tokenResult, err := apiTokenCall.Struct()
		if err != nil {
			fmt.Printf("Error getting endpoint URL: %s\n", err)
			return
		}
		token, err := tokenResult.Token()
		if err != nil {
			fmt.Printf("Error getting endpoint URL: %s\n", err)
			return
		}
		url, err := tokenResult.EndpointUrl()
		if err != nil {
			fmt.Printf("Error getting endpoint URL: %s\n", err)
			return
		}
		// this convention is respected in Sync
		pubURL := url + "#" + token
		fmt.Println("Got Pub URL %s", pubURL)
		self := feed.Pub{URL: pubURL, LastUpdated: 0, LastChecked: 0}
		db.PutSelf(&self)
	}
}
