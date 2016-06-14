package mark

import (
	"crypto/rsa"
	"encoding/json"
	"errors"

	"github.com/square/go-jose"
)

// Op is an arbitrary operation
type Op struct {
	Op       string
	OpNum    int
	FeedHash string
	KeyID    string
	Body     interface{}
}

// Feed is a sequence of operations
type Feed struct {
	Ops []Op `json:"ops"`
}

// DeclareKey returns an Op that sets the key for a feed
func DeclareKey(key *rsa.PublicKey) (*Op, error) {
	jwt, err := AsJWK(key)
	if err != nil {
		return nil, err
	}
	return &Op{Op: "declare-key", Body: jwt}, nil

}

// FromBytes inflates a Feed object from binary
func FromBytes(key *rsa.PublicKey, bytes []byte) (*Feed, error) {
	var feed Feed
	if len(bytes) == 0 {

		// bootstrap a new feed
		var ops []Op
		declareKeyOp, err := DeclareKey(key)
		if err != nil {
			return nil, err
		}
		declareKeyOp.OpNum = 0
		feed = Feed{Ops: ops}
		feed.Ops = append(feed.Ops, *declareKeyOp)
		return &feed, nil
	}
	err := json.Unmarshal(bytes, &feed)
	if err != nil {
		return nil, err
	}

	// TODO - verify feed
	return &feed, nil
}

// Append adds an Op to the end of a feed
func (feed *Feed) Append(op Op) error {
	// op.FeedHash = feed.FeedHash()
	op.OpNum = feed.Ops[len(feed.Ops)-1].OpNum + 1
	keyID, err := feed.GetKeyID()
	if err != nil {
		return err
	}
	op.KeyID = keyID
	feed.Ops = append(feed.Ops, op)
	return nil
}

// GetKeyID returns the currently delcared public key
func (feed *Feed) GetKeyID() (string, error) {
	for i := len(feed.Ops) - 1; i >= 0; i-- {
		op := feed.Ops[i]
		if op.Op == "declare-key" {
			return op.Body.(*jose.JsonWebKey).KeyID, nil
		}
	}
	return "", errors.New("Feed had no declared key")
}
