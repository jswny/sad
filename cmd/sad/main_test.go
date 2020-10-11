package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jswny/sad"
	main "github.com/jswny/sad/cmd/sad"
	testutils "github.com/jswny/sad/internal"
)

func TestGetAllOptionSources(t *testing.T) {
	expectedOpts := testutils.GetTestOpts()
	stringExpectedOpts := testutils.StringOptions{}
	stringExpectedOpts.FromOptions(&expectedOpts)

	program := "sad"

	args := buildArgs(&stringExpectedOpts)

	stringExpectedOpts.SetEnv()

	expectedOptsData, err := json.Marshal(expectedOpts)

	if err != nil {
		stringExpectedOpts.UnsetEnv()
		t.Fatalf("Error marshaling expected options: %s", err)
	}

	tempFile, err := ioutil.TempFile(".", ".sad.json.test.")

	if err != nil {
		stringExpectedOpts.UnsetEnv()
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if err := ioutil.WriteFile(tempFile.Name(), expectedOptsData, 0644); err != nil {
		stringExpectedOpts.UnsetEnv()
		t.Fatalf("Error writing to temp file: %s", err)
	}

	commandLineOpts, environmentOpts, configOpts, commandLineOutput, err := main.GetAllOptionSources(program, args, tempFile.Name())

	stringExpectedOpts.UnsetEnv()

	if err != nil {
		t.Errorf("Error getting all option sources: %s", err)
	}

	if commandLineOutput != "" {
		t.Errorf("Expected empty output but got: %s", commandLineOutput)
	}

	testutils.CompareOpts(expectedOpts, *commandLineOpts, t)
	testutils.CompareOpts(expectedOpts, *environmentOpts, t)
	testutils.CompareOpts(expectedOpts, *configOpts, t)
}

func TestMergeOptionsHierarchy(t *testing.T) {
	commandLineOpts := testutils.GetTestOpts()
	environmentOpts := testutils.GetTestOpts()
	configOpts := testutils.GetTestOpts()

	commandLineOpts.Username = ""
	commandLineOpts.RootDir = ""

	environmentOpts.RootDir = ""

	expectedOpts := sad.Options{}
	data, err := json.Marshal(commandLineOpts)
	if err != nil {
		t.Fatalf("Error marshaling command line options: %s", err)
	}

	err = json.Unmarshal(data, &expectedOpts)
	if err != nil {
		t.Fatalf("Error unmarshaling command line options: %s", err)
	}

	expectedOpts.Username = environmentOpts.Username
	expectedOpts.RootDir = configOpts.RootDir

	main.MergeOptionsHierarchy(&commandLineOpts, &environmentOpts, &configOpts)

	testutils.CompareOpts(expectedOpts, commandLineOpts, t)
}

func TestParseFlags(t *testing.T) {
	testOpts := testutils.GetTestOpts()
	stringTestOpts := testutils.StringOptions{}
	stringTestOpts.FromOptions(&testOpts)

	program := "sad"

	args := buildArgs(&stringTestOpts)

	opts, output, err := main.ParseFlags(program, args)
	if err != nil {
		t.Fatalf("Error parsing flags: %s", err)
	}

	if output != "" {
		t.Errorf("Expected empty output but got: %s", output)
	}

	testutils.CompareOpts(testOpts, *opts, t)
}

func TestParseFlagsEmptyValues(t *testing.T) {
	testOpts := sad.Options{}

	program := "sad"
	var args []string

	opts, output, err := main.ParseFlags(program, args)
	if err != nil {
		t.Fatalf("Error parsing flags: %s", err)
	}

	if output != "" {
		t.Fatalf("Expected empty output but got: %s", output)
	}

	testutils.CompareOpts(testOpts, *opts, t)
}

func buildArgs(stringOpts *testutils.StringOptions) []string {
	args := []string{
		"-name",
		stringOpts.Name,
		"-server",
		stringOpts.Server,
		"-username",
		stringOpts.Username,
		"-root-dir",
		stringOpts.RootDir,
		"-private-key",
		stringOpts.PrivateKey,
		"-channel",
		stringOpts.Channel,
		"-env-vars",
		stringOpts.EnvVars,
		"-debug",
		stringOpts.Debug,
	}

	return args
}
