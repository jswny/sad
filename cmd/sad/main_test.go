package main_test

import (
	"encoding/json"
	"testing"

	"github.com/jswny/sad"
	main "github.com/jswny/sad/cmd/sad"
	testutils "github.com/jswny/sad/internal"
)

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
		t.Fatalf("Expected empty output but got: %s", output)
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
