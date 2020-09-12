package sad_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	testutils "github.com/jswny/sad/internal"

	"github.com/jswny/sad"
)

func TestMergeNoEmptyValues(t *testing.T) {
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

func TestMergeEmptyValues(t *testing.T) {
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

func TestMergeSomeEmptyValues(t *testing.T) {
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

func TestMergeDefaults(t *testing.T) {
	expectedOpts := testutils.GetTestOpts()
	opts := sad.Options{}

	expectedOpts.Channel = ""
	expectedOpts.Path = ""

	err := testutils.CloneOptions(&expectedOpts, &opts)
	if err != nil {
		t.Fatalf("Error cloning expected options into options: %s", err)
	}

	opts.MergeDefaults()

	expectedOpts.Channel = "beta"
	expectedOpts.Path = "."

	testutils.CompareOpts(expectedOpts, opts, t)
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

	stringTestOpts.SetEnv()

	opts := sad.Options{}
	err := opts.GetEnv()

	stringTestOpts.UnsetEnv()

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
