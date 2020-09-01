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
	"strconv"
	"strings"
)

// EnvVarPrefix represents the prefix that all environment variables should have to be read properly.
var EnvVarPrefix = "SAD_"

// RSAPrivateKey wraps an RSA private key and supports conversion to/from JSON.
type RSAPrivateKey struct {
	PrivateKey *rsa.PrivateKey
}

// Options for deployment.
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

// Merge merges the other options into the existing options
// When both fields are populated, the field from the existing options is kept.
func (o *Options) Merge(other *Options) {
	if o.Server == nil {
		o.Server = other.Server
	}

	if o.Username == "" {
		o.Username = other.Username
	}

	if o.RootDir == "" {
		o.RootDir = other.RootDir
	}

	if o.PrivateKey.PrivateKey == nil {
		o.PrivateKey = other.PrivateKey
	}

	if o.Channel == "" {
		o.Channel = other.Channel
	}

	if o.Path == "" {
		o.Path = other.Path
	}

	if len(o.EnvVars) == 0 {
		o.EnvVars = other.EnvVars
	}

	if !o.Debug {
		o.Debug = other.Debug
	}
}

// FromStrings converts strings into options.
func (o *Options) FromStrings(server string, username string, rootDir string, privateKey string, channel string, path string, envVars string, debug string) error {
	if server != "" {
		o.Server = net.ParseIP(server)
	}

	o.Username = username
	o.RootDir = rootDir

	if privateKey != "" {
		rsaPrivateKey := RSAPrivateKey{}
		err := rsaPrivateKey.ParseBase64PEMString(privateKey)
		if err != nil {
			return err
		}
		o.PrivateKey = rsaPrivateKey
	}

	o.Channel = channel
	o.Path = path

	if envVars != "" {
		envVarsArr := strings.Split(envVars, ",")
		o.EnvVars = envVarsArr
	}

	if debug != "" {
		debugBool, err := strconv.ParseBool(debug)
		if err != nil {
			return err
		}

		o.Debug = debugBool
	}

	return nil
}

// GetJSON parses options from a JSON file.
func (o *Options) GetJSON(filename string) error {
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

// GetEnv parses options from environment variables.
// All variables should be prefixed and they should correspond to the available options with underscores separating words such as "PRIVATE_KEY".
// The private key should be a base64 encoded string.
// The environment variables should be a comma-separated string.
func (o *Options) GetEnv() error {
	prefix := EnvVarPrefix

	server := os.Getenv(prefix + "SERVER")
	username := os.Getenv(prefix + "USERNAME")
	rootDir := os.Getenv(prefix + "ROOT_DIR")
	privateKey := os.Getenv(prefix + "PRIVATE_KEY")
	channel := os.Getenv(prefix + "CHANNEL")
	path := os.Getenv(prefix + "PATH")
	envVars := os.Getenv(prefix + "ENV_VARS")
	debug := os.Getenv(prefix + "DEBUG")

	err := o.FromStrings(server, username, rootDir, privateKey, channel, path, envVars, debug)

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
