package feed

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"
)

const failOut = 10

// Pub represents a URL-addressable node
type Pub struct {
	URL         string `json:"url"`
	LastChecked int64  `json:"last_checked"`
	LastUpdated int64  `json:"last_updated"`
	Failures    int    `json:"failures"`
}

// Head is the length of a feed
type Head struct {
	ID  string `json:"id"`
	Len int    `json:"len"`
}

// Announcement is a pub and one or more heads for that pub
type Announcement struct {
	Pub   Pub    `json:"pub"`
	Heads []Head `json:"heads"`
}

// ShouldUpdate says whether this pub should be checked
// uses exponential backoff
func (p *Pub) ShouldUpdate() bool {
	if p.Failures > failOut {
		return false
	}
	now := time.Now().Unix()
	uncertainDuration := now - p.LastChecked
	staleDuration := p.LastChecked - p.LastUpdated
	return uncertainDuration >= staleDuration
}

// URLHash is a short way to identify pubs
func (p *Pub) URLHash() []byte {
	sha := sha256.Sum256([]byte(p.URL))
	out := make([]byte, base64.RawURLEncoding.EncodedLen(len(sha)))
	base64.RawURLEncoding.Encode(out, sha[:])
	return out
}

// GetHeads issues a request to a pub to fetch the head of each feed it has
func (p *Pub) GetHeads() ([]Head, error) {
	u, err := url.Parse(p.URL)
	u.Path = path.Join(u.Path, ProtocolRoot, HeadsPath)
	s := u.String()

	r, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var heads []Head
	err = json.NewDecoder(r.Body).Decode(&heads)
	return heads, err
}

// GetPubs issues a request to load the pubs another pub knows about
func (p *Pub) GetPubs() ([]Pub, error) {
	u, err := url.Parse(p.URL)
	u.Path = path.Join(u.Path, ProtocolRoot, PubsPath)
	s := u.String()

	r, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var pubs []Pub
	err = json.NewDecoder(r.Body).Decode(&pubs)
	return pubs, err
}

// GetFeed issues a request to load a specific feed from the pub
func (p *Pub) GetFeed(feedID string) (*SignedFeed, error) {
	u, err := url.Parse(p.URL)
	u.Path = path.Join(u.Path, ProtocolRoot, FeedPath, feedID)
	s := u.String()

	r, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var sf SignedFeed

	err = json.NewDecoder(r.Body).Decode(&sf)
	return &sf, err
}

// Announce posts an announcement to a feed
func (p *Pub) Announce(a *Announcement) error {
	u, err := url.Parse(p.URL)
	u.Path = path.Join(u.Path, ProtocolRoot, AnnouncePath)
	s := u.String()
	buf, err := json.Marshal(a)
	if err != nil {
		return err
	}

	_, err = http.Post(s, "application/json", bytes.NewBuffer(buf))
	return err
}
