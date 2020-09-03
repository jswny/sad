// Package testutils provides internal shared testing utilities.
package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	mathrand "math/rand"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/jswny/sad"
)

// StringOptions represents all options as strings.
type StringOptions struct {
	Server     string
	Username   string
	RootDir    string
	PrivateKey string
	Channel    string
	Path       string
	EnvVars    string
	Debug      string
}

// FromOptions converts options into string options.
func (stringOpts *StringOptions) FromOptions(opts *sad.Options) {
	stringOpts.Server = opts.Server.String()
	stringOpts.Username = opts.Username
	stringOpts.RootDir = opts.RootDir
	stringOpts.PrivateKey = opts.PrivateKey.ToBase64PEMString()
	stringOpts.Channel = opts.Channel
	stringOpts.Path = opts.Path
	stringOpts.EnvVars = strings.Join(opts.EnvVars, ",")
	stringOpts.Debug = strconv.FormatBool(opts.Debug)
}

// GetTestOpts retrieves a set of random options for testing.
func GetTestOpts() sad.Options {
	rsaPrivateKey := GenerateRSAPrivateKey()

	randSize := 5

	testOpts := sad.Options{
		Server:     net.ParseIP("1.2.3.4"),
		Username:   randString(randSize),
		RootDir:    randString(randSize),
		PrivateKey: rsaPrivateKey,
		Channel:    randString(randSize),
		Path:       randString(randSize),
		EnvVars: []string{
			randString(randSize),
			randString(randSize),
		},
		Debug: true,
	}

	return testOpts
}

// GenerateRSAPrivateKey generates a random RSA private key.
func GenerateRSAPrivateKey() sad.RSAPrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	rsaPrivateKey := sad.RSAPrivateKey{
		PrivateKey: privateKey,
	}
	return rsaPrivateKey
}

// CompareOpts compares two sets of options in a test environment.
func CompareOpts(expectedOpts sad.Options, actualOpts sad.Options, t *testing.T) {
	if !actualOpts.Server.Equal(expectedOpts.Server) {
		t.Errorf("Expected server IP %s but got %s", expectedOpts.Server, actualOpts.Server)
	}

	if actualOpts.Username != expectedOpts.Username {
		t.Errorf("Expected username %s but got %s", expectedOpts.Username, actualOpts.Username)
	}

	if actualOpts.RootDir != expectedOpts.RootDir {
		t.Errorf("Expected root directory %s but got %s", expectedOpts.RootDir, actualOpts.RootDir)
	}

	if (expectedOpts.PrivateKey.PrivateKey != nil && actualOpts.PrivateKey.PrivateKey != nil) && !expectedOpts.PrivateKey.PrivateKey.Equal(actualOpts.PrivateKey.PrivateKey) {
		t.Errorf("Expected equal private keys but they were not equal")
	}

	if actualOpts.Channel != expectedOpts.Channel {
		t.Errorf("Expected channel %s but got %s", expectedOpts.Channel, actualOpts.Channel)
	}

	if actualOpts.Path != expectedOpts.Path {
		t.Errorf("Expected path %s but got %s", expectedOpts.Path, actualOpts.Path)
	}

	if !testEqualSlices(actualOpts.EnvVars, expectedOpts.EnvVars) {
		t.Errorf("Expected environment variables %s but got %s", expectedOpts.EnvVars, actualOpts.EnvVars)
	}

	if actualOpts.Debug != expectedOpts.Debug {
		t.Errorf("Expected debug %t but got %t", expectedOpts.Debug, actualOpts.Debug)
	}
}

// CloneOptions clones options into other options.
// The options to clone into should ideally be empty.
func CloneOptions(optionsToClone *sad.Options, optionsToCloneInto *sad.Options) error {
	data, err := json.Marshal(optionsToClone)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &optionsToCloneInto)
	if err != nil {
		return err
	}

	return nil
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

func testEqualSlices(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
