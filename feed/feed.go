package feed

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

// Encode turns a feed into a signed feed; can only do this if you have the key to sign it
func (c *Coder) Encode(f *Feed, key *rsa.PrivateKey) (SignedFeed, error) {
	var sf SignedFeed

	for _, op := range f.Ops {
		s, err := op.ToJWS(key)
		if err != nil {
			return nil, err
		}
		sf = append(sf, s)
	}
	return sf, nil
}

// Url-safe base64 encode that strips padding
// Copied from square jose
func base64URLEncode(data []byte) string {
	var result = base64.URLEncoding.EncodeToString(data)
	return strings.TrimRight(result, "=")
}

// Url-safe base64 decoder that adds padding
// Copied from square jose
func base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	return base64.URLEncoding.DecodeString(data)
}

func contentHash(s string) string {
	sha := sha256.Sum256([]byte(s))
	out := make([]byte, base64.RawURLEncoding.EncodedLen(len(sha)))
	base64.RawURLEncoding.Encode(out, sha[:])
	return string(out)
}

// Decode turns a SignedFeed into a feed, verifying it along the way
func (c *Coder) Decode(sf SignedFeed) (*Feed, error) {
	key, err := sf.CurrentKey()
	if err != nil {
		return nil, err
	}

	var f Feed
	for i, s := range sf {
		opJws, err := jose.ParseSigned(s)
		if err != nil {
			return nil, err
		}
		opBytes, err := opJws.Verify(key)
		if err != nil {
			return nil, err
		}
		var op Op
		err = json.Unmarshal(opBytes, &op)
		if err != nil {
			return nil, err
		}

		if i > 0 {
			prev := sf[i-1]
			if op.FeedHash != contentHash(prev) {
				fmt.Printf("Verification failed for op %s\n", op)
			}
		}

		if op.OpNum != i {
			fmt.Printf("OpNum bad for %s\n", op)
		}

		op.DecodeBody(c.registry)
		f.Ops = append(f.Ops, op)
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

// ToJWS turns an op into it's JWS representation
func (op *Op) ToJWS(key *rsa.PrivateKey) (string, error) {
	signer, err := jose.NewSigner(jose.RS256, key)
	if err != nil {
		return "", err
	}
	payload, err := op.MarshalJSON()
	if err != nil {
		return "", err
	}
	jws, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}
	return jws.CompactSerialize()
}

// ContentHash retuns a sha256 of the JWS representation of an op
func (op *Op) ContentHash(key *rsa.PrivateKey) (string, error) {
	s, err := op.ToJWS(key)
	if err != nil {
		return "", err
	}
	sha := sha256.Sum256([]byte(s))
	out := make([]byte, base64.RawURLEncoding.EncodedLen(len(sha)))
	base64.RawURLEncoding.Encode(out, sha[:])
	return string(out), nil
}

// Feed is a sequence of operations
type Feed struct {
	Ops []Op
}

// SignedFeed is a feed in compact JWS serialization format
type SignedFeed []string

// CurrentKey returns the first declared key
func (sf SignedFeed) CurrentKey() (*jose.JsonWebKey, error) {
	dkJWS := sf[0]
	parts := strings.Split(dkJWS, ".")
	dkPayload := parts[1]
	dkBytes, err := base64URLDecode(dkPayload)
	if err != nil {
		return nil, err
	}
	var dkOp Op
	err = json.Unmarshal(dkBytes, &dkOp)
	if err != nil {
		return nil, err
	}
	var jwk jose.JsonWebKey
	err = json.Unmarshal(dkOp.RawBody, &jwk)
	return &jwk, err
}

// Fingerprint returns a fingerprint of a pub key
func (sf SignedFeed) Fingerprint() (string, error) {
	jwk, err := sf.CurrentKey()
	if err != nil {
		return "", err
	}
	return fingerprint(jwk)
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

// FeedHash returns a content hash of the current latest op
func (feed *Feed) FeedHash(key *rsa.PrivateKey) (string, error) {
	prev := feed.Ops[len(feed.Ops)-1]
	return prev.ContentHash(key)
}

// Append adds an Op to the end of a feed
func (feed *Feed) Append(op Op, key *rsa.PrivateKey) error {
	fh, err := feed.FeedHash(key)
	if err != nil {
		return err
	}
	op.FeedHash = fh
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
		}
	}
	return nil, errors.New("Feed had no declared key")
}

// Len returns the length of a feed
func (feed *Feed) Len() int {
	return len(feed.Ops)
}

func fingerprint(jwk *jose.JsonWebKey) (string, error) {
	thumbprint, err := jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", err
	}
	out := make([]byte, base64.RawURLEncoding.EncodedLen(len(thumbprint)))
	base64.RawURLEncoding.Encode(out, thumbprint)
	return string(out), nil
}

// Fingerprint returns a fingerprint of a pub key
func (feed *Feed) Fingerprint() (string, error) {
	jwk, err := feed.CurrentKey()
	if err != nil {
		return "", err
	}
	return fingerprint(jwk)
}

// Fingerprint returns a fingerprint of a pub key
func Fingerprint(key *rsa.PublicKey) (string, error) {
	jwk, err := AsJWK(key)
	if err != nil {
		return "", err
	}
	return fingerprint(jwk)
}
