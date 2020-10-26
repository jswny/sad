// Package testutils provides internal shared testing utilities.
package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	mathrand "math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jswny/sad"
)

// StringOptions represents all options as strings.
type StringOptions struct {
	Repository  string
	Server      string
	Username    string
	RootDir     string
	PrivateKey  string
	Channel     string
	Path        string
	EnvVars     string
	Debug       string
	ImageDigest string
}

// FromOptions converts options into string options.
func (stringOpts *StringOptions) FromOptions(opts *sad.Options) {
	stringOpts.Repository = opts.Repository
	stringOpts.Server = opts.Server.String()
	stringOpts.Username = opts.Username
	stringOpts.RootDir = opts.RootDir
	stringOpts.PrivateKey = opts.PrivateKey.ToBase64PEMString()
	stringOpts.Channel = opts.Channel
	stringOpts.EnvVars = strings.Join(opts.EnvVars, ",")
	stringOpts.Debug = strconv.FormatBool(opts.Debug)
	stringOpts.ImageDigest = opts.ImageDigest
}

// SetEnv sets environment variables for all string options.
// UnsetEnv should be called after.
func (stringOpts *StringOptions) SetEnv() {
	prefix := sad.EnvVarPrefix

	setEnvFromPrefixPostfix(prefix, "REPOSITORY", stringOpts.Repository)
	setEnvFromPrefixPostfix(prefix, "SERVER", stringOpts.Server)
	setEnvFromPrefixPostfix(prefix, "USERNAME", stringOpts.Username)
	setEnvFromPrefixPostfix(prefix, "ROOT_DIR", stringOpts.RootDir)
	setEnvFromPrefixPostfix(prefix, "PRIVATE_KEY", stringOpts.PrivateKey)
	setEnvFromPrefixPostfix(prefix, "CHANNEL", stringOpts.Channel)
	setEnvFromPrefixPostfix(prefix, "ENV_VARS", stringOpts.EnvVars)
	setEnvFromPrefixPostfix(prefix, "DEBUG", stringOpts.Debug)
	setEnvFromPrefixPostfix(prefix, "IMAGE_DIGEST", stringOpts.ImageDigest)
}

// UnsetEnv sets environment variables for all string options.
// Should be called after SetEnv.
func (stringOpts *StringOptions) UnsetEnv() {
	prefix := sad.EnvVarPrefix
	var envVarPostfix string

	envVarPostfix = "REPOSITORY"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "SERVER"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "USERNAME"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ROOT_DIR"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PRIVATE_KEY"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "CHANNEL"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PATH"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ENV_VARS"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "DEBUG"
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "IMAGE_DIGEST"
	defer os.Unsetenv(prefix + envVarPostfix)
}

// GetTestOpts retrieves a set of random options for testing.
func GetTestOpts() sad.Options {
	rsaPrivateKey := GenerateRSAPrivateKey()

	randSize := 5

	testOpts := sad.Options{
		Repository: randString(randSize),
		Server:     net.ParseIP("1.2.3.4"),
		Username:   randString(randSize),
		RootDir:    randString(randSize),
		PrivateKey: rsaPrivateKey,
		Channel:    randString(randSize),
		EnvVars: []string{
			randString(randSize),
			randString(randSize),
		},
		Debug:       true,
		ImageDigest: randString(randSize),
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
	CompareStrings("repository", expectedOpts.Repository, actualOpts.Repository, t)

	if !actualOpts.Server.Equal(expectedOpts.Server) {
		t.Errorf("Expected server IP %s but got %s", expectedOpts.Server, actualOpts.Server)
	}

	CompareStrings("username", expectedOpts.Username, actualOpts.Username, t)

	CompareStrings("root directory", expectedOpts.RootDir, actualOpts.RootDir, t)

	if (expectedOpts.PrivateKey.PrivateKey != nil && actualOpts.PrivateKey.PrivateKey != nil) && !expectedOpts.PrivateKey.PrivateKey.Equal(actualOpts.PrivateKey.PrivateKey) {
		t.Errorf("Expected equal private keys but they were not equal")
	}

	CompareStrings("channel", expectedOpts.Channel, actualOpts.Channel, t)

	compareSlices("environment variables", expectedOpts.EnvVars, actualOpts.EnvVars, t)

	if expectedOpts.Debug != actualOpts.Debug {
		t.Errorf("Expected debug %t but got %t", expectedOpts.Debug, actualOpts.Debug)
	}

	CompareStrings("image digest", expectedOpts.ImageDigest, actualOpts.ImageDigest, t)
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

// CompareStrings compares two strings in a test and nicely handles empty strings.
func CompareStrings(name, expected, actual string, t *testing.T) {
	if expected != actual {
		empty := "<empty>"

		if expected == "" {
			expected = empty
		}

		if actual == "" {
			actual = empty
		}

		t.Errorf("Expected %s %s but got %s", name, expected, actual)
	}
}

// CompareReaderLines compares the contents of a reader to a slice of lines.
func CompareReaderLines(name string, expectedLines []string, actualReader io.Reader, t *testing.T) {
	actual := ReadFromReader(name, actualReader, t)

	expectedNumLines := len(expectedLines)

	actualLines := strings.Split(actual, "\n")

	actualNumLines := len(actualLines)

	var strippedActualLines []string

	for _, actualLine := range actualLines {
		if actualLine != "" {
			strippedActualLines = append(strippedActualLines, actualLine)
		}
	}

	strippedActualNumLines := len(strippedActualLines)

	if strippedActualNumLines != expectedNumLines {
		t.Errorf("Expected %d non-empty lines in reader contents for %s but got %d out of a total of %d lines", expectedNumLines, name, strippedActualNumLines, actualNumLines)
	}

	for _, expectedLine := range expectedLines {
		if !strings.Contains(actual, expectedLine) {
			t.Errorf("Expected line \"%s\" in reader contents for %s but got:\n%s", expectedLine, name, actual)
		}
	}
}

// ReadFromReader reads a reader into a string or errors out fatally.
func ReadFromReader(name string, reader io.Reader, t *testing.T) string {
	buffer := new(strings.Builder)
	_, err := io.Copy(buffer, reader)

	if err != nil {
		t.Fatalf("Error reading from reader for %s: %s", name, err)
	}

	return buffer.String()
}

// SetEnvVarsWithPrefix sets the environment variables defined by the specified options.
// Each variable is prefixed by the default prefix.
func SetEnvVarsWithPrefix(opts *sad.Options, content string) {
	for _, variableName := range opts.EnvVars {
		os.Setenv(sad.EnvVarPrefix+variableName, content)
	}
}

// UnsetEnvVarsWithPrefix unsets the environment variables defined by the specified options.
// Each variable is prefixed by the default prefix.
func UnsetEnvVarsWithPrefix(opts *sad.Options) {
	for _, variableName := range opts.EnvVars {
		os.Unsetenv(sad.EnvVarPrefix + variableName)
	}
}

func compareSlices(name string, expected, actual []string, t *testing.T) {
	equal := true

	if (expected == nil) != (actual == nil) {
		equal = false
	}

	if len(expected) != len(actual) {
		equal = false
	}

	for i := range expected {
		if expected[i] != actual[i] {
			equal = false
		}
	}

	if !equal {
		t.Errorf("Expected %s %s but got %s", name, expected, actual)
	}
}

func randString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathrand.Intn(len(letters))]
	}
	return string(b)
}

func setEnvFromPrefixPostfix(prefix, postfix, value string) {
	os.Setenv(prefix+postfix, value)
}
