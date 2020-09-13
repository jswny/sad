package sad

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"

	"golang.org/x/crypto/ssh"
)

// RSAPrivateKey wraps an RSA private key and supports conversion to/from JSON.
type RSAPrivateKey struct {
	PrivateKey *rsa.PrivateKey
}

// MarshalJSON marshals an RSA private key into valid JSON.
// The key is marshalled into a base64 encoded PEM block string.
func (k RSAPrivateKey) MarshalJSON() ([]byte, error) {
	encoded := k.ToBase64PEMString()

	marshaledData, err := json.Marshal(encoded)

	if err != nil {
		return nil, errors.New("Failed to marshal encoded pem data to JSON")
	}

	return marshaledData, nil
}

// UnmarshalJSON unmarshals JSON into an RSA private key.
// The key should be a base64 encoded PEM block string.
func (k *RSAPrivateKey) UnmarshalJSON(data []byte) error {
	var unmarshaled string
	err := json.Unmarshal(data, &unmarshaled)

	if err != nil {
		return err
	}

	err = k.ParseBase64PEMString(unmarshaled)

	if err != nil {
		return err
	}

	return nil
}

// ToBase64PEMString converts an RSA private key into a base 64 encoded PEM block string.
func (k *RSAPrivateKey) ToBase64PEMString() string {
	var data []byte
	if k.PrivateKey != nil {
		data = x509.MarshalPKCS1PrivateKey(k.PrivateKey)
	}

	pemBlock := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: data,
		},
	)

	encoded := base64.StdEncoding.EncodeToString(pemBlock)

	return encoded
}

// ParseBase64PEMString parses a base64 encoded PEM block string into an RSA private key.
func (k *RSAPrivateKey) ParseBase64PEMString(str string) error {
	decoded, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		return err
	}

	block, _ := pem.Decode(decoded)

	if block == nil {
		return errors.New("Failed to parse PEM block containing RSA private key")
	}

	var privateKey *rsa.PrivateKey
	if len(block.Bytes) > 0 {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	if err != nil {
		return err
	}

	k.PrivateKey = privateKey
	return nil
}

// ToSSHAuthMethod converts an RSA private key into an SSH AuthMethod.
func (k *RSAPrivateKey) ToSSHAuthMethod() (ssh.AuthMethod, error) {
	signer, err := ssh.NewSignerFromKey(k.PrivateKey)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
