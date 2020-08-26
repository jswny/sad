package sad

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net"
	"os"
)

// RSAPrivateKey wraps an RSA private key and supports conversion to/from JSON
type RSAPrivateKey struct {
	PrivateKey *rsa.PrivateKey
}

// Options for deployment
type Options struct {
	Server     net.IP
	Username   string
	RootDir    string
	PrivateKey RSAPrivateKey
	Channel    string
	Path       string
	EnvVars    []string
	Debug      bool
}

// MarshalJSON marshals an RSA private key into valid JSON
func (k RSAPrivateKey) MarshalJSON() ([]byte, error) {
	data := x509.MarshalPKCS1PrivateKey(k.PrivateKey)
	pemBlock := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: data,
		},
	)

	marshaledData, err := json.Marshal(pemBlock)

	if err != nil {
		return nil, errors.New("Failed to marshal encoded pem data to JSON")
	}

	return marshaledData, nil
}

// UnmarshalJSON unmarshals JSON into an RSA private key
func (k *RSAPrivateKey) UnmarshalJSON(data []byte) error {
	var unmarshaled string
	err := json.Unmarshal(data, &unmarshaled)
	if err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(unmarshaled)

	if err != nil {
		return err
	}

	block, _ := pem.Decode(decoded)

	if block == nil {
		return errors.New("Failed to parse PEM block containing RSA private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return err
	}

	k.PrivateKey = privateKey
	return err
}

// Hello says hello to the world
func Hello() string {
	return "Hello, world."
}

// Get parses options from a file
func (o *Options) Get(filename string) error {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, o)
}
