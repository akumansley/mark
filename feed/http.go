package feed

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Getter is an interface for a get-only http service
type Getter interface {
	Get(string) (*http.Response, error)
}

var once sync.Once
var instance Getter

// Initialize sets up the http service with a getter implementation
// First caller wins
func Initialize(getter Getter) {
	once.Do(func() {
		instance = getter
	})
}

// Get executes an HTTP get
// Calling Get before Initialize is an error
func Get(requestedURL string) (*http.Response, error) {
	if instance == nil {
		panic("Calling Get before initialization")
	}

	// Handle special sandstorm URL format
	splot := strings.Split(requestedURL, "#")
	if len(splot) == 2 {
		u, err := url.Parse(splot[0])
		if err != nil {
			return nil, err
		}
		u.User = url.UserPassword("anonymous-user", splot[1])
		// HACK: if we know the node is a sandstorm node,
		// we have to drop "sync" from the path as sandstorm will add it automatically
		// as a defense mechanism
		p := u.EscapedPath()
		if strings.HasPrefix(p, "/sync") {
			p = p[len("/sync"):]
		}
		u.Path = p

		requestedURL = u.String()
	}
	return instance.Get(requestedURL)
}
