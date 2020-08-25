package sad_test

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
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

func TestOptionsGet(t *testing.T) {
	privateKeyData := "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"

	privateKey, err := ParseRsaPrivateKeyFromPemStr(privateKeyData)

	if err != nil {
		t.Fatalf("Error parsing test private key: %s", err)
	}

	testOpts := sad.Options{
		Server:     net.ParseIP("1.2.3.4"),
		Username:   "user1",
		RootDir:    "/srv",
		PrivateKey: privateKey,
		Channel:    "beta",
		Path:       "/app",
		EnvVars:    []string{"foo", "bar"},
		Debug:      true,
	}

	testOptsData, err := json.Marshal(testOpts)

	if err != nil {
		t.Fatalf("Error marshaling test options struct: %s", err)
	}

	tempFile, err := ioutil.TempFile("", ".sad.json.test.")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if err := ioutil.WriteFile(tempFile.Name(), testOptsData, 0644); err != nil {
		t.Fatalf("Error writing to temp file: %s", err)
	}

	opts := sad.Options{}

	if err := opts.Get(tempFile.Name()); err != nil {
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

	actualPrivateKeyString := ExportRsaPrivateKeyAsPemStr(opts.PrivateKey)
	expectedPrivateKeyString := ExportRsaPrivateKeyAsPemStr(testOpts.PrivateKey)
	if actualPrivateKeyString != expectedPrivateKeyString {
		t.Errorf("Expected private key %s but got %s", expectedPrivateKeyString, actualPrivateKeyString)
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

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("Failed to parse PEM block containing SSH private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem)
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
