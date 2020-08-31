package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	testutils "github.com/jswny/sad/internal"
)

var binName = "sad"

func TestMain(m *testing.M) {
	fmt.Println("Building CLI...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build CLI %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)

	os.Exit(result)
}

func TestCLI(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("Run", func(t *testing.T) {
		cmd := exec.Command(cmdPath)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("RunWithArgs", func(t *testing.T) {
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

		cmd := exec.Command(
			cmdPath,
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
		)

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
}
