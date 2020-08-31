// Package testutils provides internal shared testing utilities.
package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"net"
	"strconv"
	"strings"

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

	testOpts := sad.Options{
		Server:     net.ParseIP("1.2.3.4"),
		Username:   "user1",
		RootDir:    "/srv",
		PrivateKey: rsaPrivateKey,
		Channel:    "beta",
		Path:       "/app",
		EnvVars:    []string{"foo", "bar"},
		Debug:      true,
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
