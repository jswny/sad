package sad_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jswny/sad"
	testutils "github.com/jswny/sad/internal"
)

func TestFindFilePathRecursive(t *testing.T) {
	tempDirName := "dir.test"
	tempFileName := "file.test"

	tempDirPath, err := ioutil.TempDir("", tempDirName)

	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	defer os.RemoveAll(tempDirPath)

	tempFile, err := ioutil.TempFile(tempDirPath, tempFileName)

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	generatedTempFileName := filepath.Base(tempFile.Name())

	stringPathSeparator := string(os.PathSeparator)
	splitTempDirPath := strings.Split(tempDirPath, stringPathSeparator)
	splitTempDirPathLen := len(splitTempDirPath)
	mainTempDir := strings.Join(splitTempDirPath[:splitTempDirPathLen-1], stringPathSeparator)

	path, err := sad.FindFilePathRecursive(mainTempDir, generatedTempFileName)

	if err != nil {
		t.Fatalf("Error finding recursive file path: %s", err)
	}

	expected := tempFile.Name()
	actual := path

	testutils.CompareStrings(expected, actual, "file path", t)
}

func TestFindFilePathRecursiveCWD(t *testing.T) {
	fileName := "file.test"

	tempFile, err := ioutil.TempFile(".", fileName)

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	defer os.Remove(tempFile.Name())

	tempFileName := filepath.Base(tempFile.Name())

	path, err := sad.FindFilePathRecursive(".", tempFileName)

	if err != nil {
		t.Fatalf("Error finding recursive file path: %s", err)
	}

	expected := tempFile.Name()
	actual := path

	testutils.CompareStrings(expected, actual, "file path", t)
}

func TestFindFilePathRecursiveNotFound(t *testing.T) {
	cwdPath, err := os.Getwd()

	if err != nil {
		t.Fatalf("Error getting current working directory: %s", err)
	}

	fileName := "1234556789"
	path, err := sad.FindFilePathRecursive(cwdPath, fileName)

	if err == nil {
		t.Fatalf("Expected error but got nil!")
	}

	testutils.CompareStrings(sad.FindFilePathRecursiveFileNotFoundErrorMessage, err.Error(), "error", t)

	expected := ""
	actual := path

	testutils.CompareStrings(expected, actual, "file path", t)
}
