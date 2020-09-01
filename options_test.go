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

func TestRSAPrivateKeyMarshalJSONNil(t *testing.T) {
	rsaPrivateKey := sad.RSAPrivateKey{}

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

func TestRSAPrivateKeyUnmarshalJSONNil(t *testing.T) {
	rsaPrivateKey := sad.RSAPrivateKey{}

	firstKeyData, _ := rsaPrivateKey.MarshalJSON()

	rsaPrivateKey2 := sad.RSAPrivateKey{}

	err := rsaPrivateKey2.UnmarshalJSON(firstKeyData)

	if err != nil {
		t.Fatalf("Error unmarshaling RSA private key: %s", err)
	}

	if rsaPrivateKey.PrivateKey != rsaPrivateKey2.PrivateKey {
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

	testutils.CompareOpts(testOpts, opts, t)
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

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsGetJSONEmtpyValues(t *testing.T) {
	testOpts := sad.Options{}
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

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsGetEnv(t *testing.T) {
	testOpts := testutils.GetTestOpts()
	stringTestOpts := testutils.StringOptions{}
	stringTestOpts.FromOptions(&testOpts)

	prefix := sad.EnvVarPrefix

	var envVarPostfix string

	envVarPostfix = "SERVER"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.Server)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "USERNAME"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.Username)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ROOT_DIR"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.RootDir)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PRIVATE_KEY"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.PrivateKey)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "CHANNEL"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.Channel)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PATH"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.Path)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ENV_VARS"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.EnvVars)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "DEBUG"
	os.Setenv(prefix+envVarPostfix, stringTestOpts.Debug)
	defer os.Unsetenv(prefix + envVarPostfix)

	opts := sad.Options{}
	err := opts.GetEnv()

	if err != nil {
		t.Fatalf("Error getting options from environment: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsGetEnvEmptyValues(t *testing.T) {
	testOpts := sad.Options{}

	opts := sad.Options{}
	err := opts.GetEnv()

	if err != nil {
		t.Fatalf("Error getting options from environment: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}
