package feed

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Pub struct {
	URL         string `json:"url"`
	LastChecked int64  `json:"last_checked"`
	LastUpdated int64  `json:"last_updated"`
}

type Head struct {
	ID  string `json:"id"`
	Len int    `json:"len"`
}

func (p *Pub) ShouldUpdate() bool {
	now := time.Now().Unix()
	uncertainDuration := now - p.LastChecked
	staleDuration := p.LastChecked - p.LastUpdated
	return uncertainDuration >= staleDuration
}

func (p *Pub) URLHash() []byte {
	sha := sha256.Sum256([]byte(p.URL))
	out := make([]byte, base64.RawURLEncoding.EncodedLen(len(sha)))
	base64.RawURLEncoding.Encode(out, sha[:])
	return out
}

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

func (p *Pub) GetFeed(feedID string) (*Feed, error) {
	u, err := url.Parse(p.URL)
	u.Path = path.Join(u.Path, ProtocolRoot, FeedPath, feedID)
	s := u.String()

	r, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var feed Feed

	// TODO this ain't gonna work
	err = json.NewDecoder(r.Body).Decode(&feed)
	return &feed, err
}
