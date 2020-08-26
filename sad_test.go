package sad_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/jswny/sad"
)

func TestHello(t *testing.T) {
	s := sad.Hello()
	expected := "Hello, world."
	if s != expected {
		t.Errorf("Expected %q, got %q instead", expected, s)
	}
}

func TestRSAPrivateKeyMarshalJSON(t *testing.T) {
	privateKey := generatePrivateKey()

	rsaPrivateKey := sad.RSAPrivateKey{
		PrivateKey: privateKey,
	}

	data, err := rsaPrivateKey.MarshalJSON()

	if err != nil {
		t.Fatalf("Error marshaling first RSA private key: %s", err)
	}

	if !json.Valid(data) {
		t.Errorf("RSA private key marshal to JSON did not produce valid JSON. Got %s", data)
	}
}

func TestRSAPrivateKeyUnmarshalJSON(t *testing.T) {
	privateKey := generatePrivateKey()

	rsaPrivateKey := sad.RSAPrivateKey{
		PrivateKey: privateKey,
	}

	firstKeyData, _ := rsaPrivateKey.MarshalJSON()

	rsaPrivateKey2 := sad.RSAPrivateKey{}

	err := rsaPrivateKey2.UnmarshalJSON(firstKeyData)

	if err != nil {
		t.Fatalf("Error unmarshaling RSA private key: %s", err)
	}

	if !rsaPrivateKey.PrivateKey.Equal(rsaPrivateKey2.PrivateKey) {
		t.Errorf("Expected marshaled and unmarshaled private keys to be equal, but they were not equal")
	}
}

func TestOptionsGet(t *testing.T) {
	privateKey := generatePrivateKey()
	rsaPrivateKey := sad.RSAPrivateKey{PrivateKey: privateKey}

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

	testOptsData, err := json.Marshal(testOpts)

	if err != nil {
		t.Fatalf("Error marshaling test options: %s", err)
	}

	tempFile, err := ioutil.TempFile(".", ".sad.json.test.")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if err := ioutil.WriteFile(tempFile.Name(), testOptsData, 0644); err != nil {
		t.Fatalf("Error writing to temp file: %s", err)
	}

	opts := sad.Options{}

	if err := opts.GetJSON(tempFile.Name()); err != nil {
		t.Fatalf("Error getting options from file: %s", err)
	}

	if !opts.Server.Equal(testOpts.Server) {
		t.Errorf("Expected server IP %s but got %s", testOpts.Server, opts.Server)
	}

	if opts.Username != testOpts.Username {
		t.Errorf("Expected username %s but got %s", testOpts.Username, opts.Username)
	}

	if opts.RootDir != testOpts.RootDir {
		t.Errorf("Expected root directory %s but got %s", testOpts.RootDir, opts.RootDir)
	}

	if !testOpts.PrivateKey.PrivateKey.Equal(opts.PrivateKey.PrivateKey) {
		t.Errorf("Expected equal private keys but they were not equal")
	}

	if opts.Channel != testOpts.Channel {
		t.Errorf("Expected channel %s but got %s", testOpts.Channel, opts.Channel)
	}

	if opts.Path != testOpts.Path {
		t.Errorf("Expected path %s but got %s", testOpts.Path, opts.Path)
	}

	if !testEqualSlices(opts.EnvVars, testOpts.EnvVars) {
		t.Errorf("Expected environment variables %s but got %s", testOpts.EnvVars, opts.EnvVars)
	}

	if opts.Debug != testOpts.Debug {
		t.Errorf("Expected debug %t but got %t", testOpts.Debug, opts.Debug)
	}
}

func generatePrivateKey() *rsa.PrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privateKey
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
