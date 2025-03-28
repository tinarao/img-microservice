package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type KeysPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func NewKeysPair(email string) (*KeysPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	publicKey := &privateKey.PublicKey

	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(email), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt public key: %v", err)
	}

	publicKeyStr := base64.StdEncoding.EncodeToString(encrypted)
	privateKeyStr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))

	return &KeysPair{
		PublicKey:  publicKeyStr,
		PrivateKey: privateKeyStr,
	}, nil
}

func (kp *KeysPair) Unmarshal() (content *string, err error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(kp.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	encrypted, err := base64.StdEncoding.DecodeString(kp.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt public key: %v", err)
	}

	c := string(decryptedBytes)
	return &c, nil
}
