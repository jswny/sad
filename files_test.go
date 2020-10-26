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

	testutils.CompareStrings("file path", expected, actual, t)
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

	testutils.CompareStrings("file path", expected, actual, t)
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

	testutils.CompareStrings("error", sad.FindFilePathRecursiveFileNotFoundErrorMessage, err.Error(), t)

	expected := ""
	actual := path

	testutils.CompareStrings("file path", expected, actual, t)
}

func TestGetEntitesForDeployment(t *testing.T) {
	dirName := "dir.test"

	tempDirPath, err := ioutil.TempDir("", dirName)

	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	defer os.RemoveAll(tempDirPath)

	fileNames := []string{
		sad.LocalDockerComposeFileName,
	}

	content := []byte("test")

	for _, fileName := range fileNames {
		filePath := filepath.Join(tempDirPath, fileName)

		if err := ioutil.WriteFile(filePath, content, 0755); err != nil {
			t.Fatalf("Error writing to temp file \"%s\", %s", filePath, err)
		}
	}

	variables := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	envVarNames := make([]string, len(variables))

	i := 0
	for variableName := range variables {
		envVarNames[i] = variableName
		i++
	}

	opts := &sad.Options{
		Repository:  "user/repo",
		Channel:     "beta",
		EnvVars:     envVarNames,
		ImageDigest: "abc123",
	}

	variableContent := "test"

	for _, variableName := range opts.EnvVars {
		os.Setenv(variableName, variableContent)
		defer os.Unsetenv(variableName)
	}

	readerMap, files, err := sad.GetEntitiesForDeployment(tempDirPath, opts)

	if err != nil {
		t.Fatalf("Error getting files for deployment: %s", err)
	}

	expected := 1 + len(fileNames)
	actual := len(readerMap)

	if actual != expected {
		t.Errorf("Returned %d readers, expected %d", actual, expected)
	}

	expected = 2
	actual = len(readerMap)

	if actual != expected {
		t.Errorf("Returned %d readers, expected %d", actual, expected)
	}

	composeReader := readerMap[sad.RemoteDockerComposeFileName]

	name := "Docker Compose file"

	data := testutils.ReadFromReader(name, composeReader, t)

	testutils.CompareStrings(name, string(content), data, t)

	reader := readerMap[sad.RemoteDotEnvFileName]

	expectedContent := []string{
		"foo=test",
		"baz=test",
		"IMAGE=user/repo@abc123",
		"CONTAINER_NAME=user-repo-beta",
	}

	testutils.CompareReaderLines(".env file", expectedContent, reader, t)

	expected = len(fileNames)
	actual = len(files)

	if actual != expected {
		t.Errorf("Returned %d files, expected %d", actual, expected)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file.Name())

		if err != nil {
			t.Fatalf("Error reading returned from file: %s", err)
		}

		if string(content) != string(data) {
			t.Errorf("Expected returned file content %s but got %s", content, data)
		}
	}
}

func TestGenerateDotEnvFile(t *testing.T) {
	variables := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	reader := sad.GenerateDotEnvFile(variables)

	expected := []string{
		"foo=bar\n",
		"baz=qux\n",
	}

	testutils.CompareReaderLines(".env file", expected, reader, t)
}

func TestFilesToFileNameReaderMap(t *testing.T) {
	var tempFiles []*os.File

	fileName := "file.test"
	content := []byte("test")
	numFiles := 3

	for i := 0; i < numFiles; i++ {
		tempFile, err := ioutil.TempFile(".", fileName)

		if err != nil {
			t.Fatalf("Error creating temp file: %s", err)
		}

		defer os.Remove(tempFile.Name())

		filePath := tempFile.Name()

		if err := ioutil.WriteFile(filePath, content, 0755); err != nil {
			t.Fatalf("Error writing to temp file \"%s\", %s", filePath, err)
		}

		tempFiles = append(tempFiles, tempFile)
	}

	actual := sad.FilesToFileNameReaderMap(tempFiles)
	numActual := len(actual)

	expected := tempFiles
	numExpected := len(tempFiles)

	if len(actual) != len(tempFiles) {
		t.Errorf("Expected %d items in map, got %d", numExpected, numActual)
	}

	for _, tempFile := range expected {
		fileName := filepath.Base(tempFile.Name())

		reader := actual[fileName]
		if reader == nil {
			t.Errorf("Reader for file %s was nil", fileName)
		}

		actualContent := testutils.ReadFromReader(fileName, reader, t)

		if actualContent != string(content) {
			t.Errorf("Expected content %s for file %s, got %s", content, fileName, actualContent)
		}
	}
}
