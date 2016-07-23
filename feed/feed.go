package feed

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"encoding/base64"
	"errors"

	"github.com/square/go-jose"
)

// Converter deserializes an op's body
type Converter func([]byte) (interface{}, error)

// Coder can marshal and unmarshal a feed
type Coder struct {
	registry map[string]Converter
}

// NewCoder returns a new coder
func NewCoder() *Coder {
	return &Coder{registry: make(map[string]Converter)}
}

// RegisterOp tells a coder how to interpret an op type's body
func (c *Coder) RegisterOp(name string, cf Converter) {
	c.registry[name] = cf
}

// Encode turns a feed into bytes
func (c *Coder) Encode(f *Feed) ([]byte, error) {
	return json.Marshal(f)
}

// Decode turns bytes into a feed
func (c *Coder) Decode(bytes []byte) (*Feed, error) {
	var f Feed

	json.Unmarshal(bytes, &f)
	for i := range f.Ops {
		f.Ops[i].DecodeBody(c.registry)
	}

	return &f, nil
}

// Op is an arbitrary operation
type Op struct {
	Op       string
	OpNum    int
	FeedHash string
	RawBody  json.RawMessage `json:"Body"`
	Body     interface{}     `json:"-"`
}

// MarshalJSON mostly just copies the Body to the RawBody
func (op *Op) MarshalJSON() ([]byte, error) {
	var rawBody []byte
	var err error
	if op.Body != nil {
		rawBody, err = json.Marshal(op.Body)
		if err != nil {
			return nil, err
		}
	} else {
		rawBody = op.RawBody
	}
	obj := make(map[string]interface{})
	obj["OpNum"] = op.OpNum
	obj["Op"] = op.Op
	obj["FeedHash"] = op.FeedHash
	raw := json.RawMessage(rawBody)
	obj["Body"] = &raw
	return json.Marshal(obj)
}

// DecodeBody loads an op's body from json
func (op *Op) DecodeBody(registry map[string]Converter) error {
	body, err := registry[op.Op](op.RawBody)
	op.Body = body
	return err
}

// Feed is a sequence of operations
type Feed struct {
	Ops []Op
}

// DeclareKey returns an Op that sets the key for a feed
func DeclareKey(key *rsa.PublicKey) (*Op, error) {
	jwk, err := AsJWK(key)
	if err != nil {
		return nil, err
	}
	return &Op{Op: "declare-key", Body: jwk}, nil
}

// New bootstraps a feed
func New(key *rsa.PrivateKey) (*Feed, error) {
	var ops []Op
	declareKeyOp, err := DeclareKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	declareKeyOp.OpNum = 0
	feed := Feed{Ops: ops}
	feed.Ops = append(feed.Ops, *declareKeyOp)
	return &feed, nil
}

// FromBytes inflates a Feed object from binary
func FromBytes(bytes []byte) (*Feed, error) {
	var feed Feed
	if len(bytes) == 0 {

		return nil, errors.New("bytes is empty")
	}
	err := json.Unmarshal(bytes, &feed)
	if err != nil {
		return nil, err
	}

	// TODO - verify feed
	return &feed, nil
}

// ToBytes serializes a feed
func (feed *Feed) ToBytes(key *rsa.PrivateKey) ([]byte, error) {
	signer, err := jose.NewSigner(jose.RS256, key)
	if err != nil {
		return nil, err
	}
	var bytesList [][]byte

	for _, op := range feed.Ops {
		payload, err := json.Marshal(op)
		if err != nil {
			return nil, err
		}
		jws, err := signer.Sign(payload)
		if err != nil {
			return nil, err
		}
		s, err := jws.CompactSerialize()
		if err != nil {
			return nil, err
		}
		bytes := []byte(s)
		bytesList = append(bytesList, bytes)
	}
	return json.Marshal(bytesList)
}

// Append adds an Op to the end of a feed
func (feed *Feed) Append(op Op) error {
	// op.FeedHash = feed.FeedHash()
	op.OpNum = feed.Ops[len(feed.Ops)-1].OpNum + 1
	feed.Ops = append(feed.Ops, op)
	return nil
}

// CurrentKey returns the currently delcared public key
func (feed *Feed) CurrentKey() (*jose.JsonWebKey, error) {
	for i := len(feed.Ops) - 1; i >= 0; i-- {
		op := feed.Ops[i]
		if op.Op == "declare-key" {
			// this case indicates that i should do something better at unmarshal time
			if op.Body != nil {
				jwk := op.Body.(*jose.JsonWebKey)
				return jwk, nil
			}
			var jwk jose.JsonWebKey
			err := json.Unmarshal(op.RawBody, &jwk)
			if err != nil {
				return nil, err
			}
			return &jwk, nil
		}
	}
	return nil, errors.New("Feed had no declared key")
}

func fingerprint(jwk *jose.JsonWebKey) ([]byte, error) {
	thumbprint, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	out := make([]byte, base64.RawURLEncoding.EncodedLen(len(thumbprint)))
	base64.RawURLEncoding.Encode(out, thumbprint)
	return out, nil
}

// Fingerprint returns a fingerprint of a pub key
func (feed *Feed) Fingerprint() ([]byte, error) {
	jwk, err := feed.CurrentKey()
	if err != nil {
		return nil, err
	}
	return fingerprint(jwk)
}

// Fingerprint returns a fingerprint of a pub key
func Fingerprint(key *rsa.PublicKey) ([]byte, error) {
	jwk, err := AsJWK(key)
	if err != nil {
		return nil, err
	}
	return fingerprint(jwk)
}
