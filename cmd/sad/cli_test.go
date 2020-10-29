package main_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	args := []string{"-h"}
	cmd, out, err := generateCmd(args)

	if err != nil {
		t.Fatalf("Error generating command to execute: %s", err)
	}

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() != 2 {
			t.Fatalf("Error executing command: %s", err)
		}
	}

	if !strings.Contains(out.String(), "Usage of") {
		t.Errorf("Help message was not printed")
	}
}

func TestCLIInvalidFlag(t *testing.T) {
	args := []string{"-invalid"}
	cmd, out, err := generateCmd(args)

	if err != nil {
		t.Fatalf("Error generating command to execute: %s", err)
	}

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() != 1 {
			t.Fatalf("Error executing command: %s", err)
		}
	}

	if !strings.Contains(out.String(), "flag provided but not defined:") {
		t.Errorf("Flag error was not printed")
	}

	if !strings.Contains(out.String(), "Usage of") {
		t.Errorf("Help message was not printed")
	}
}

func generateCmd(args []string) (*exec.Cmd, *bytes.Buffer, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	cmdPath := filepath.Join(dir, binName)

	var out bytes.Buffer
	cmd := exec.Command(cmdPath, args...)
	cmd.Stdout = &out

	return cmd, &out, nil
}
