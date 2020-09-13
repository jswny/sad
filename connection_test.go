package sad_test

import (
	"fmt"
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

func TestGetSSHClientConfig(t *testing.T) {
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

func TestGetSCPClient(t *testing.T) {
	opts := testutils.GetTestOpts()

	clientConfig, err := sad.GetSSHClientConfig(&opts)

	if err != nil {
		t.Fatalf("Error getting SSH client config: %s", err)
	}

	scpClient, err := sad.GetSCPClient(&opts, clientConfig)

	if err != nil {
		t.Fatalf("Error getting SCP client: %s", err)
	}

	expected := fmt.Sprintf("%s:%d", opts.Server.String(), 22)
	actual := scpClient.Host

	testutils.CompareStrings(expected, actual, "SCP client host", t)
}
