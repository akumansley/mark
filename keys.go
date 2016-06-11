package mark

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"path"
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

	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)

	pubPemData := pem.EncodeToMemory(&pem.Block{
		Type:  pubKeyType,
		Bytes: pubBytes,
	})

	privPemData := pem.EncodeToMemory(&pem.Block{
		Type:  privKeyType,
		Bytes: privBytes,
	})

	publicKeyPath := path.Join(markDir, publicKeyFilename)
	ioutil.WriteFile(publicKeyPath, pubPemData, 0644)

	privateKeyPath := path.Join(markDir, privateKeyFilename)
	ioutil.WriteFile(privateKeyPath, privPemData, 0600)

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

	privKeyBlock, _ := pem.Decode(privKeyData)
	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	pubKeyBlock, _ := pem.Decode(pubKeyData)
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	privKey.PublicKey = *pubKey.(*rsa.PublicKey)

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
