package sad_test

import (
	"io"
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

func TestGetFilesForDeployment(t *testing.T) {
	dirName := "dir.test"

	tempDirPath, err := ioutil.TempDir("", dirName)

	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	defer os.RemoveAll(tempDirPath)

	fileNames := []string{
		sad.DockerComposeFileName,
	}

	content := []byte("test")

	for _, fileName := range fileNames {
		filePath := filepath.Join(tempDirPath, fileName)

		if err := ioutil.WriteFile(filePath, content, 0755); err != nil {
			t.Fatalf("Error writing to temp file \"%s\", %s", filePath, err)
		}
	}

	files, err := sad.GetFilesForDeployment(tempDirPath)

	if err != nil {
		t.Fatalf("Error getting files for deployment: %s", err)
	}

	expected := len(fileNames)
	actual := len(files)

	if actual != expected {
		t.Errorf("Getting files for deployment returned %d files, expected %d", actual, expected)
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file.Name())

		if err != nil {
			t.Fatalf("Error reading from deployment file: %s", err)
		}

		if string(content) != string(data) {
			t.Errorf("Expected file content %s but got %s", content, data)
		}
	}
}

func TestGenerateDotEnvFile(t *testing.T) {
	variables := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	reader := sad.GenerateDotEnvFile(variables)
	builder := new(strings.Builder)
	_, err := io.Copy(builder, reader)

	if err != nil {
		t.Fatalf("Error copying reader to string builder: %s", err)
	}

	actual := builder.String()
	expected := []string{
		"foo=bar\n",
		"baz=qux\n",
	}

	for _, expectedLine := range expected {
		if !strings.Contains(actual, expectedLine) {
			t.Errorf("Expected line %s in .env contents but got:\n%s", expectedLine, actual)
		}
	}
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

		builder := new(strings.Builder)
		_, err := io.Copy(builder, reader)

		if err != nil {
			t.Fatalf("Error copying reader for file %s to string builder: %s", fileName, err)
		}

		actualContent := builder.String()

		if actualContent != string(content) {
			t.Errorf("Expected content %s for file %s, got %s", content, fileName, actualContent)
		}
	}
}
