package feed

import (
	"crypto/sha256"
	"encoding/base64"
	"time"
)

type Pub struct {
	URL         string `json:"url"`
	LastChecked int64  `json:"last_checked"`
	LastUpdated int64  `json:"last_updated"`
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
