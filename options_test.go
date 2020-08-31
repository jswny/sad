package sad_test

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	testutils "github.com/jswny/sad/internal"

	"github.com/jswny/sad"
)

func TestRSAPrivateKeyMarshalJSON(t *testing.T) {
	rsaPrivateKey := testutils.GenerateRSAPrivateKey()

	data, err := rsaPrivateKey.MarshalJSON()

	if err != nil {
		t.Fatalf("Error marshaling first RSA private key: %s", err)
	}

	if !json.Valid(data) {
		t.Errorf("RSA private key marshal to JSON did not produce valid JSON. Got %s", data)
	}
}

func TestRSAPrivateKeyUnmarshalJSON(t *testing.T) {
	rsaPrivateKey := testutils.GenerateRSAPrivateKey()

	firstKeyData, _ := rsaPrivateKey.MarshalJSON()

	rsaPrivateKey2 := sad.RSAPrivateKey{}

	err := rsaPrivateKey2.UnmarshalJSON(firstKeyData)

	if err != nil {
		t.Fatalf("Error unmarshaling RSA private key: %s", err)
	}

	if !rsaPrivateKey.PrivateKey.Equal(rsaPrivateKey2.PrivateKey) {
		t.Errorf("Expected marshaled and unmarshaled private keys to be equal, but they were not")
	}
}

func TestToBase64PEMString(t *testing.T) {
	rsaPrivateKey := testutils.GenerateRSAPrivateKey()
	encoded := rsaPrivateKey.ToBase64PEMString()

	_, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Errorf("PEM block string was not valid base64 encoding")
	}
}

func TestParseBase64PEMString(t *testing.T) {
	testRSAPrivateKey := testutils.GenerateRSAPrivateKey()
	encoded := testRSAPrivateKey.ToBase64PEMString()

	rsaPrivateKey := sad.RSAPrivateKey{}
	err := rsaPrivateKey.ParseBase64PEMString(encoded)

	if err != nil {
		t.Fatalf("Failed to parse base64 PEM string into an RSA private key")
	}

	if !testRSAPrivateKey.PrivateKey.Equal(rsaPrivateKey.PrivateKey) {
		t.Errorf("Expected base64 PEM block encoded and decoded private keys to be equal, but they were not")
	}
}

func TestOptionsFromStrings(t *testing.T) {
	testOpts := testutils.GetTestOpts()
	stringTestOpts := testutils.StringOptions{}
	stringTestOpts.FromOptions(&testOpts)

	server := stringTestOpts.Server
	username := stringTestOpts.Username
	rootDir := stringTestOpts.RootDir
	privateKey := stringTestOpts.PrivateKey
	channel := stringTestOpts.Channel
	path := stringTestOpts.Path
	envVars := stringTestOpts.EnvVars
	debug := stringTestOpts.Debug

	opts := sad.Options{}
	err := opts.FromStrings(server, username, rootDir, privateKey, channel, path, envVars, debug)
	if err != nil {
		t.Fatalf("Error getting options from test options strings: %s", err)
	}

	compareOpts(testOpts, opts, t)
}

func TestOptionsGetJSON(t *testing.T) {
	testOpts := testutils.GetTestOpts()
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
	testOpts := testutils.GetTestOpts()
	stringTestOpts := testutils.StringOptions{}
	stringTestOpts.FromOptions(&testOpts)

	prefix := sad.EnvVarPrefix

	os.Setenv(prefix+"SERVER", stringTestOpts.Server)
	defer os.Unsetenv("SERVER")

	os.Setenv(prefix+"USERNAME", stringTestOpts.Username)
	defer os.Unsetenv("USERNAME")

	os.Setenv(prefix+"ROOT_DIR", stringTestOpts.RootDir)
	defer os.Unsetenv("ROOT_DIR")

	os.Setenv(prefix+"PRIVATE_KEY", stringTestOpts.PrivateKey)
	defer os.Unsetenv("PRIVATE_KEY")

	os.Setenv(prefix+"CHANNEL", stringTestOpts.Channel)
	defer os.Unsetenv("CHANNEL")

	os.Setenv(prefix+"PATH", stringTestOpts.Path)
	defer os.Unsetenv("PATH")

	os.Setenv(prefix+"ENV_VARS", stringTestOpts.EnvVars)
	defer os.Unsetenv("ENV_VARS")

	os.Setenv(prefix+"DEBUG", stringTestOpts.Debug)
	defer os.Unsetenv("DEBUG")

	opts := sad.Options{}
	err := opts.GetEnv()

	if err != nil {
		t.Fatalf("Error getting options from environment: %s", err)
	}

	compareOpts(testOpts, opts, t)
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
