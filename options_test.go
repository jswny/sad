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

	repository := stringTestOpts.Repository
	server := stringTestOpts.Server
	username := stringTestOpts.Username
	rootDir := stringTestOpts.RootDir
	privateKey := stringTestOpts.PrivateKey
	channel := stringTestOpts.Channel
	envVars := stringTestOpts.EnvVars
	debug := stringTestOpts.Debug
	imageDigest := stringTestOpts.ImageDigest

	opts := sad.Options{}
	err := opts.FromStrings(repository, server, username, rootDir, privateKey, channel, envVars, debug, imageDigest)
	if err != nil {
		t.Fatalf("Error getting options from test options strings: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsFromJSON(t *testing.T) {
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

	if err := opts.FromJSON(tempFile.Name()); err != nil {
		t.Fatalf("Error getting options from file: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsFromJSONEmtpyValues(t *testing.T) {
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

	if err := opts.FromJSON(tempFile.Name()); err != nil {
		t.Fatalf("Error getting options from file: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsFromEnv(t *testing.T) {
	testOpts := testutils.GetTestOpts()
	stringTestOpts := testutils.StringOptions{}
	stringTestOpts.FromOptions(&testOpts)

	stringTestOpts.SetEnv()
	defer stringTestOpts.UnsetEnv()

	opts := sad.Options{}
	err := opts.FromEnv()

	if err != nil {
		t.Fatalf("Error getting options from environment: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}

func TestOptionsFromEnvEmptyValues(t *testing.T) {
	testOpts := sad.Options{}

	opts := sad.Options{}
	err := opts.FromEnv()

	if err != nil {
		t.Fatalf("Error getting options from environment: %s", err)
	}

	testutils.CompareOpts(testOpts, opts, t)
}

func TestGetDeploymentName(t *testing.T) {
	opts := sad.Options{
		Repository: "user/foo",
		Channel:    "beta/123",
	}

	deploymentName, err := opts.GetDeploymentName()

	if err != nil {
		t.Fatalf("Error getting deployment name: %s", err)
	}

	expected := "user-foo-beta-123"

	testutils.CompareStrings("full name", expected, deploymentName, t)
}

func TestGetImageSpecifier(t *testing.T) {
	opts := sad.Options{
		Repository:  "user/foo",
		ImageDigest: "sha256:abc123",
	}

	deploymentName := opts.GetImageSpecifier()

	expected := "user/foo@sha256:abc123"

	testutils.CompareStrings("full name", expected, deploymentName, t)
}

func TestGetEnvValues(t *testing.T) {
	opts := sad.Options{
		EnvVars: []string{
			"foo",
			"bar",
		},
	}

	content := "test"

	testutils.SetEnvVarsWithPrefix(&opts, content)
	defer testutils.UnsetEnvVarsWithPrefix(&opts)

	envMap, err := opts.GetEnvValues()

	if err != nil {
		t.Fatalf("Error getting referenced environment variables: %s", err)
	}

	for _, variableName := range opts.EnvVars {
		variableValue := envMap[variableName]

		name := fmt.Sprintf("environment variable %s value", variableName)

		testutils.CompareStrings(name, content, variableValue, t)
	}
}

func TestGetEnvValuesBlank(t *testing.T) {
	opts := sad.Options{
		EnvVars: []string{
			"foo",
			"bar",
		},
	}

	content := "test"

	testutils.SetEnvVarsWithPrefix(&opts, content)

	toUnset := opts.EnvVars[0]
	os.Unsetenv(sad.EnvVarPrefix + toUnset)
	defer testutils.UnsetEnvVarsWithPrefix(&opts)

	envMap, err := opts.GetEnvValues()

	if err == nil {
		t.Errorf("Expected error getting referenced environment variables but got: %s", err)
	}

	containsVarName := strings.Contains(err.Error(), toUnset)
	if !containsVarName {
		t.Errorf("Expected error getting referenced environment variable %s but got: %s", toUnset, err)
	}

	if envMap != nil {
		t.Errorf("Expected nil returned environment variables but got: %s", envMap)
	}
}
