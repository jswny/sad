package sad_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/jswny/sad"
)

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

func TestOptionsGetJSON(t *testing.T) {
	testOpts := getTestOpts()
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

	compareOpts(testOpts, opts, t)
}

func TestOptionsGetEnv(t *testing.T) {
	testOpts := getTestOpts()

	os.Setenv("SERVER", testOpts.Server.String())
	defer os.Unsetenv("SERVER")

	os.Setenv("USERNAME", testOpts.Username)
	defer os.Unsetenv("USERNAME")

	os.Setenv("ROOT_DIR", testOpts.RootDir)
	defer os.Unsetenv("ROOT_DIR")

	encoded := testOpts.PrivateKey.ToBase64PEMData()

	os.Setenv("PRIVATE_KEY", string(encoded))
	defer os.Unsetenv("PRIVATE_KEY")

	os.Setenv("CHANNEL", testOpts.Channel)
	defer os.Unsetenv("CHANNEL")

	os.Setenv("PATH", testOpts.Path)
	defer os.Unsetenv("PATH")

	envVars := strings.Join(testOpts.EnvVars, ",")
	os.Setenv("ENV_VARS", envVars)
	defer os.Unsetenv("ENV_VARS")

	debug := strconv.FormatBool(testOpts.Debug)
	os.Setenv("DEBUG", debug)
	defer os.Unsetenv("DEBUG")

	opts := sad.Options{}
	err := opts.GetEnv()

	if err != nil {
		t.Fatalf("Error getting options from environment: %s", err)
	}
}

func compareOpts(expectedOpts sad.Options, actualOpts sad.Options, t *testing.T) {
	if !actualOpts.Server.Equal(expectedOpts.Server) {
		t.Errorf("Expected server IP %s but got %s", expectedOpts.Server, actualOpts.Server)
	}

	if actualOpts.Username != expectedOpts.Username {
		t.Errorf("Expected username %s but got %s", expectedOpts.Username, actualOpts.Username)
	}

	if actualOpts.RootDir != expectedOpts.RootDir {
		t.Errorf("Expected root directory %s but got %s", expectedOpts.RootDir, actualOpts.RootDir)
	}

	if !expectedOpts.PrivateKey.PrivateKey.Equal(actualOpts.PrivateKey.PrivateKey) {
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

func getTestOpts() sad.Options {
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

	return testOpts
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
