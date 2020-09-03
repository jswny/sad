package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
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

func TestCLIHelpMessage(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	cmd := exec.Command(cmdPath, "-h")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() != 2 {
			t.Fatal(err)
		}
	}
}
