package sad

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DockerComposeFileName is the name of the Docker Compose file to be loaded for deployment.
var DockerComposeFileName string = ".sad.docker-compose.yml"

// DotEnvFileName is the name of the .env file to be loaded for deployment.
var DotEnvFileName string = ".sad.env"

// FindFilePathRecursiveFileNotFoundErrorMessage is the string error message returned when FindFilePathRecursive cannot find the specified file.
var FindFilePathRecursiveFileNotFoundErrorMessage = "file not found"

// FindFilePathRecursive finds a file path recursively that matches the specified file name starting from the specified path.
// Returns the path of the file if it is found, otherwise returns an error.
// If the error was only that the file was not found, returns an error containing FindFilePathRecursiveFileNotFoundErrorMessage.
func FindFilePathRecursive(fromPath string, fileName string) (string, error) {
	var foundPath string

	foundErrorMessage := "file found"

	err := filepath.Walk(fromPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && fileName == info.Name() {
			foundPath = path
			return errors.New(foundErrorMessage)
		}
		return nil
	})

	if err != nil && err.Error() != foundErrorMessage {
		return "", err
	}

	if foundPath == "" {
		return "", errors.New(FindFilePathRecursiveFileNotFoundErrorMessage)
	}

	return foundPath, nil
}

// GetFilesForDeployment gets and opens the files needed for deployment by finding them recursively under the provided fromPath.
// Files: Docker Compose file (see DockerComposeFileName).
// Opens files, remember to close.
func GetFilesForDeployment(fromPath string) ([]*os.File, error) {
	var filePaths []string
	var files []*os.File

	fileNames := []string{
		DockerComposeFileName,
	}

	for _, fileName := range fileNames {
		filePath, err := FindFilePathRecursive(fromPath, fileName)

		if err != nil {
			err := fmt.Errorf("error finding file \"%s\" under path \"%s\": %s", fileName, fromPath, err)
			return nil, err
		}

		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)

		if err != nil {
			err := fmt.Errorf("error opening file for deployment from path \"%s\"", filePath)
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

// GenerateDotEnvFile generates a file as a reader which contains a properly-formatted .env file.
func GenerateDotEnvFile(variables map[string]string) io.Reader {
	var s string

	for name, value := range variables {
		s += fmt.Sprintf("%s=%s\n", name, value)
	}

	return strings.NewReader(s)
}

// FilesToFileNameReaderMap converts a slice of files into a map of the file name to a reader for the content.
func FilesToFileNameReaderMap(files []*os.File) map[string]io.Reader {
	m := make(map[string]io.Reader)

	for _, file := range files {
		fileName := filepath.Base(file.Name())
		m[fileName] = file
	}

	return m
}
