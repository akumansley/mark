package feed

import (
	"net/http"
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
func Get(url string) (*http.Response, error) {
	if instance == nil {
		panic("Calling Get before initialization")
	}
	return instance.Get(url)
}
