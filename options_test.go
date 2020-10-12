package sad_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	testutils "github.com/jswny/sad/internal"

	"github.com/jswny/sad"
)

func TestOptionsMergeNoEmptyValues(t *testing.T) {
	expectedOpts := testutils.GetTestOpts()
	optsToMerge := testutils.GetTestOpts()
	optsToMergeInto := sad.Options{}

	err := testutils.CloneOptions(&expectedOpts, &optsToMergeInto)
	if err != nil {
		t.Fatalf("Error cloning expected options into options to merge into: %s", err)
	}

	optsToMergeInto.Merge(&optsToMerge)

	testutils.CompareOpts(expectedOpts, optsToMergeInto, t)
}

func TestOptionsMergeEmptyValues(t *testing.T) {
	expectedOpts := testutils.GetTestOpts()
	optsToMerge := sad.Options{}
	optsToMergeInto := sad.Options{}

	err := testutils.CloneOptions(&expectedOpts, &optsToMergeInto)
	if err != nil {
		t.Fatalf("Error cloning expected options into options to merge into: %s", err)
	}

	optsToMergeInto.Merge(&optsToMerge)

	testutils.CompareOpts(expectedOpts, optsToMergeInto, t)
}

func TestOptionsMergeSomeEmptyValues(t *testing.T) {
	expectedOpts := testutils.GetTestOpts()
	optsToMerge := sad.Options{}
	optsToMergeInto := sad.Options{}

	expectedOpts.Username = ""
	expectedOpts.RootDir = ""

	err := testutils.CloneOptions(&expectedOpts, &optsToMergeInto)
	if err != nil {
		t.Fatalf("Error cloning expected options into options to merge into: %s", err)
	}

	optsToMergeInto.Merge(&optsToMerge)

	expectedOpts.Username = optsToMerge.Username
	expectedOpts.RootDir = optsToMerge.RootDir

	testutils.CompareOpts(expectedOpts, optsToMergeInto, t)
}

func TestOptionsMergeDefaults(t *testing.T) {
	expectedOpts := testutils.GetTestOpts()
	opts := sad.Options{}

	expectedOpts.Channel = ""

	err := testutils.CloneOptions(&expectedOpts, &opts)
	if err != nil {
		t.Fatalf("Error cloning expected options into options: %s", err)
	}

	opts.MergeDefaults()

	expectedOpts.Channel = "beta"

	testutils.CompareOpts(expectedOpts, opts, t)
}

func TestOptionsVerifyValid(t *testing.T) {
	opts := testutils.GetTestOpts()

	err := opts.Verify()

	if err != nil {
		t.Errorf("Error verifying options: %s", err)
	}
}

func TestOptionsVerifyInvalid(t *testing.T) {
	opts := testutils.GetTestOpts()
	opts.Username = ""

	err := opts.Verify()

	if err == nil {
		t.Errorf("No error verifying options")
	}

	if !strings.ContainsAny(err.Error(), "username is <empty>") {
		t.Errorf("Error message doesn't contain username error")
	}
}

func TestOptionsFromStrings(t *testing.T) {
	testOpts := testutils.GetTestOpts()
	stringTestOpts := testutils.StringOptions{}
	stringTestOpts.FromOptions(&testOpts)

	name := stringTestOpts.Name
	server := stringTestOpts.Server
	username := stringTestOpts.Username
	rootDir := stringTestOpts.RootDir
	privateKey := stringTestOpts.PrivateKey
	channel := stringTestOpts.Channel
	envVars := stringTestOpts.EnvVars
	debug := stringTestOpts.Debug

	opts := sad.Options{}
	err := opts.FromStrings(name, server, username, rootDir, privateKey, channel, envVars, debug)
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

	stringTestOpts.SetEnv()
	defer stringTestOpts.UnsetEnv()

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

func TestGetFullAppName(t *testing.T) {
	opts := sad.Options{
		Name:    "foo",
		Channel: "beta",
	}

	fullName := opts.GetFullAppName()
	expected := "foo-beta"

	testutils.CompareStrings(expected, fullName, "full name", t)
}

func TestGetEnvValues(t *testing.T) {
	opts := sad.Options{
		EnvVars: []string{
			"foo",
			"bar",
		},
	}

	content := "test"

	for _, variableName := range opts.EnvVars {
		os.Setenv(variableName, content)
		defer os.Unsetenv(variableName)
	}

	envMap := opts.GetEnvValues()

	for _, variableName := range opts.EnvVars {
		variableValue := envMap[variableName]

		name := fmt.Sprintf("environment variable %s value", variableName)

		testutils.CompareStrings(content, variableValue, name, t)
	}
}
