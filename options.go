package sad

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// EnvVarPrefix represents the prefix that all environment variables should have to be read properly
var EnvVarPrefix = "SAD_"

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
// The key is marshalled into a base64 encoded PEM block string
func (k RSAPrivateKey) MarshalJSON() ([]byte, error) {
	encoded := k.ToBase64PEMString()

	marshaledData, err := json.Marshal(encoded)

	if err != nil {
		return nil, errors.New("Failed to marshal encoded pem data to JSON")
	}

	return marshaledData, nil
}

// UnmarshalJSON unmarshals JSON into an RSA private key
// The key should be a base64 encoded PEM block string
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

// GetJSON parses options from a JSON file
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

// GetEnv parses options from environment variables
// All variables should be prefixed and they should correspond to the available options with underscores separating words such as `PRIVATE_KEY`
// The private key should be a base64 encoded string
// The environment variables should be a comma-separated string
func (o *Options) GetEnv() error {
	prefix := EnvVarPrefix

	if envVar := os.Getenv(prefix + "SERVER"); envVar != "" {
		o.Server = net.ParseIP(envVar)
	}

	if envVar := os.Getenv(prefix + "USERNAME"); envVar != "" {
		o.Username = envVar
	}

	if envVar := os.Getenv(prefix + "ROOT_DIR"); envVar != "" {
		o.RootDir = envVar
	}

	if envVar := os.Getenv(prefix + "PRIVATE_KEY"); envVar != "" {
		k := RSAPrivateKey{}
		err := k.ParseBase64PEMString(envVar)

		if err != nil {
			return err
		}

		o.PrivateKey = k
	}

	if envVar := os.Getenv(prefix + "CHANNEL"); envVar != "" {
		o.Channel = envVar
	}

	if envVar := os.Getenv(prefix + "PATH"); envVar != "" {
		o.Path = envVar
	}

	if envVar := os.Getenv(prefix + "ENV_VARS"); envVar != "" {
		envVarsArr := strings.Split(envVar, ",")
		o.EnvVars = envVarsArr
	}

	if envVar := os.Getenv(prefix + "DEBUG"); envVar != "" {
		debug, err := strconv.ParseBool(envVar)

		if err != nil {
			o.Debug = false
		}

		o.Debug = debug
	}

	return nil
}

// GetFlags parses options from command line flags
func (o *Options) GetFlags() error {
	server := flag.String("server", "", "Server to deploy to")
	username := flag.String("username", "", "User to login to on the server")
	rootDir := flag.String("root-dir", "", "Root directory to deploy to on the server")
	privateKey := flag.String("private-key", "", "Base64 encoded SSH private key to login to the user on the server")
	channel := flag.String("channel", "", "Deployment channel")
	path := flag.String("path", "", "Path to the app to be deployed relative to the current directory")
	envVars := flag.String("env-vars", "", "Local environment variables to be injected into the app deployment")
	debug := flag.Bool("debug", false, "Debug mode")

	flag.Parse()

	o.Server = net.ParseIP(*server)
	o.Username = *username
	o.RootDir = *rootDir

	rsaPrivateKey := RSAPrivateKey{}
	err := rsaPrivateKey.ParseBase64PEMString(*privateKey)
	if err != nil {
		return err
	}
	o.PrivateKey = rsaPrivateKey

	o.Channel = *channel
	o.Path = *path

	envVarsArr := strings.Split(*envVars, ",")
	o.EnvVars = envVarsArr

	o.Debug = *debug

	return nil
}

// ToBase64PEMString converts an RSA private key into a base 64 encoded PEM block string
func (k *RSAPrivateKey) ToBase64PEMString() string {
	data := x509.MarshalPKCS1PrivateKey(k.PrivateKey)
	pemBlock := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: data,
		},
	)

	encoded := base64.StdEncoding.EncodeToString(pemBlock)

	return encoded
}

// ParseBase64PEMString parses a base64 encoded PEM block string into an RSA private key
func (k *RSAPrivateKey) ParseBase64PEMString(str string) error {
	decoded, err := base64.StdEncoding.DecodeString(str)

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
	return nil
}
