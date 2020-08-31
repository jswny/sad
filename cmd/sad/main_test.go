package main_test

import (
	"testing"

	main "github.com/jswny/sad/cmd/sad"
	testutils "github.com/jswny/sad/internal"
)

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
