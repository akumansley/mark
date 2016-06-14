package mark

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"path"

	"github.com/square/go-jose"
)

const (
	privateKeyFilename = "key"
	publicKeyFilename  = "key.pub"
	pubKeyType         = "RSA PUBLIC KEY"
	privKeyType        = "RSA PRIVATE KEY"
)

// CreateKeys makes a public private keypair and saves them in markDir
func CreateKeys(markDir string) (*rsa.PrivateKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateJWK := jose.JsonWebKey{
		Key:       privKey,
		Algorithm: string(jose.RSA1_5),
	}
	thumbprint, err := privateJWK.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	privateJWK.KeyID = string(thumbprint)
	if !privateJWK.Valid() {
		return nil, errors.New("invalid private key")
	}

	publicJWK := jose.JsonWebKey{
		Key:       &privKey.PublicKey,
		Algorithm: string(jose.RSA1_5),
	}

	thumbprint, err = publicJWK.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, err
	}
	publicJWK.KeyID = string(thumbprint)

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

// Fingerprint returns a fingerprint of a pub key
// Following camlistore convention, we return a hash name followed by a hex
// encoded digest of the data
// In this case, we use only sha256
func Fingerprint(key *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	bytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	hash.Write(bytes)
	return append([]byte("sha256-"), []byte(hex.EncodeToString(hash.Sum(nil)))...), nil
}
