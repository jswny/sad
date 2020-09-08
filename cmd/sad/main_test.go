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

	server := stringExpectedOpts.Server
	username := stringExpectedOpts.Username
	rootDir := stringExpectedOpts.RootDir
	privateKey := stringExpectedOpts.PrivateKey
	channel := stringExpectedOpts.Channel
	path := stringExpectedOpts.Path
	envVars := stringExpectedOpts.EnvVars
	debug := stringExpectedOpts.Debug

	program := "sad"
	args := []string{
		"-server",
		server,
		"-username",
		username,
		"-root-dir",
		rootDir,
		"-private-key",
		privateKey,
		"-channel",
		channel,
		"-path",
		path,
		"-env-vars",
		envVars,
		"-debug",
		debug,
	}

	prefix := sad.EnvVarPrefix
	var envVarPostfix string

	envVarPostfix = "SERVER"
	os.Setenv(prefix+envVarPostfix, server)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "USERNAME"
	os.Setenv(prefix+envVarPostfix, username)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ROOT_DIR"
	os.Setenv(prefix+envVarPostfix, rootDir)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PRIVATE_KEY"
	os.Setenv(prefix+envVarPostfix, stringExpectedOpts.PrivateKey)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "CHANNEL"
	os.Setenv(prefix+envVarPostfix, stringExpectedOpts.Channel)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "PATH"
	os.Setenv(prefix+envVarPostfix, stringExpectedOpts.Path)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "ENV_VARS"
	os.Setenv(prefix+envVarPostfix, stringExpectedOpts.EnvVars)
	defer os.Unsetenv(prefix + envVarPostfix)

	envVarPostfix = "DEBUG"
	os.Setenv(prefix+envVarPostfix, stringExpectedOpts.Debug)
	defer os.Unsetenv(prefix + envVarPostfix)

	expectedOptsData, err := json.Marshal(expectedOpts)

	if err != nil {
		t.Fatalf("Error marshaling expected options: %s", err)
	}

	tempFile, err := ioutil.TempFile(".", ".sad.json.test.")

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if err := ioutil.WriteFile(tempFile.Name(), expectedOptsData, 0644); err != nil {
		t.Fatalf("Error writing to temp file: %s", err)
	}

	commandLineOpts, environmentOpts, configOpts, commandLineOutput, err := main.GetAllOptionSources(program, args, tempFile.Name())

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

	server := stringTestOpts.Server
	username := stringTestOpts.Username
	rootDir := stringTestOpts.RootDir
	privateKey := stringTestOpts.PrivateKey
	channel := stringTestOpts.Channel
	path := stringTestOpts.Path
	envVars := stringTestOpts.EnvVars
	debug := stringTestOpts.Debug

	program := "sad"
	args := []string{
		"-server",
		server,
		"-username",
		username,
		"-root-dir",
		rootDir,
		"-private-key",
		privateKey,
		"-channel",
		channel,
		"-path",
		path,
		"-env-vars",
		envVars,
		"-debug",
		debug,
	}

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
