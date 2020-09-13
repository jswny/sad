package sad_test

import (
	"testing"

	testutils "github.com/jswny/sad/internal"

	"github.com/jswny/sad"
)

func TestGetFullName(t *testing.T) {
	opts := sad.Options{
		Name:    "foo",
		Channel: "beta",
	}

	fullName := sad.GetFullName(&opts)
	expected := "foo-beta"

	testutils.CompareStrings(expected, fullName, "full name", t)
}

func TestGetClientConfig(t *testing.T) {
	opts := testutils.GetTestOpts()

	clientConfig, err := sad.GetSSHClientConfig(&opts)

	if err != nil {
		t.Fatalf("Error getting SSH client config: %s", err)
	}

	authMethodCount := len(clientConfig.Auth)
	if authMethodCount != 1 {
		t.Errorf("Expected one auth method in SSH client config, got %d", authMethodCount)
	}
}
