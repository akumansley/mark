package feed

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"path"

	"github.com/square/go-jose"
)

// TODO move this to the app package
const (
	privateKeyFilename = "key"
	publicKeyFilename  = "key.pub"
)

// CreateKeys makes a public private keypair and saves them in markDir
func CreateKeys(markDir string) (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateJWK, err := AsJWK(privKey)
	if err != nil {
		return nil, err
	}
	if !privateJWK.Valid() {
		return nil, errors.New("invalid private key")
	}

	publicJWK, err := AsJWK(&privKey.PublicKey)
	if err != nil {
		return nil, err
	}
	if !publicJWK.Valid() {
		return nil, errors.New("invalid public key")
	}

	publicKeyPath := path.Join(markDir, publicKeyFilename)
	bytes, err := publicJWK.MarshalJSON()
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile(publicKeyPath, bytes, 0644)

	privateKeyPath := path.Join(markDir, privateKeyFilename)
	bytes, err = privateJWK.MarshalJSON()
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile(privateKeyPath, bytes, 0600)

	return OpenKeys(markDir)
}

// OpenKeys reads a public/private keypair and prepares them for use
func OpenKeys(markDir string) (*rsa.PrivateKey, error) {
	privKeyPath := path.Join(markDir, privateKeyFilename)
	pubKeyPath := path.Join(markDir, publicKeyFilename)

	privKeyData, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	pubKeyData, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}

	var privJWK, pubJWK jose.JsonWebKey

	err = privJWK.UnmarshalJSON(privKeyData)
	if err != nil {
		return nil, err
	}
	err = pubJWK.UnmarshalJSON(pubKeyData)
	if err != nil {
		return nil, err
	}

	privKey := privJWK.Key.(*rsa.PrivateKey)
	privKey.PublicKey = *pubJWK.Key.(*rsa.PublicKey)

	// Calculations that speed up private key operations in the future
	privKey.Precompute()

	// Validate Private Key -- Sanity checks on the key
	if err = privKey.Validate(); err != nil {
		return nil, err
	}

	return privKey, nil
}

// AsJWK returns a JWK representation of a key
func AsJWK(key interface{}) (*jose.JsonWebKey, error) {
	JWK := jose.JsonWebKey{
		Key:       key,
		Algorithm: string(jose.RSA1_5),
	}
	thumbprint, err := JWK.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	JWK.KeyID = base64.URLEncoding.EncodeToString(thumbprint)
	return &JWK, nil
}
