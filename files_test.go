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
	dirName := "dir.test"
	fileName := "file.test"

	tempDirPath, err := ioutil.TempDir("", dirName)

	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	defer os.RemoveAll(tempDirPath)

	tempFile, err := ioutil.TempFile(tempDirPath, fileName)

	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	tempFileName := filepath.Base(tempFile.Name())

	stringPathSeparator := string(os.PathSeparator)
	splitTempDirPath := strings.Split(tempDirPath, stringPathSeparator)
	splitTempDirPathLen := len(splitTempDirPath)
	mainTempDir := strings.Join(splitTempDirPath[:splitTempDirPathLen-1], stringPathSeparator)

	path, err := sad.FindFilePathRecursive(mainTempDir, tempFileName)

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
