package sad_test

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/jswny/sad"
	testutils "github.com/jswny/sad/internal"
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

	if err := rsaPrivateKey.PrivateKey.Validate(); err != nil {
		t.Errorf("Unmarshalled private key was not valid")
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

func TestRSAPrivateKeyToBase64PEMString(t *testing.T) {
	rsaPrivateKey := testutils.GenerateRSAPrivateKey()
	encoded := rsaPrivateKey.ToBase64PEMString()

	_, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		t.Errorf("PEM block string was not valid base64 encoding")
	}
}

func TestRSAPrivateKeyParseBase64PEMString(t *testing.T) {
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

func TestRSAPrivateKeyToSSHAuthMethod(t *testing.T) {
	testRSAPrivateKey := testutils.GenerateRSAPrivateKey()
	_, err := testRSAPrivateKey.ToSSHAuthMethod()

	if err != nil {
		t.Errorf("Error converting RSA private key to SSH auth method")
	}
}
