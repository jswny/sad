package sad

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalDockerComposeFileName is the name of the local Docker Compose file to loade for deployment.
var LocalDockerComposeFileName string = ".sad.docker-compose.yml"

// RemoteDockerComposeFileName is the name of the remote Docker Compose file to send to the server.
var RemoteDockerComposeFileName string = "docker-compose.yml"

// RemoteDotEnvFileName is the name of the remote .env file to send to the server.
var RemoteDotEnvFileName string = ".env"

// ConfigFileName is the name of the configuration file to pull options from.
var ConfigFileName string = ".sad.json"

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

// GetEntitiesForDeployment gets (and opens if necessary) the entities needed for deployment.
// Files are locaed by finding them recursively under the provided path .
// Files: Docker Compose file (see DockerComposeFileName).
// Other: generated .env file.
// Files are only returned so they can be closed by the caller.
func GetEntitiesForDeployment(fromPath string, opts *Options) (map[string]io.Reader, []*os.File, error) {
	files, err := getFilesForDeployment(fromPath)

	if err != nil {
		return nil, files, fmt.Errorf("error getting files for deployment: %s", err)
	}

	readerMap := FilesToFileNameReaderMap(files)
	readerMap[RemoteDockerComposeFileName] = readerMap[LocalDockerComposeFileName]
	delete(readerMap, LocalDockerComposeFileName)

	env := opts.GetEnvValues()
	containerName := opts.GetFullAppName()
	env["CONTAINER_NAME"] = containerName

	readerMap[RemoteDotEnvFileName] = GenerateDotEnvFile(env)

	return readerMap, files, nil
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

func getFilesForDeployment(fromPath string) ([]*os.File, error) {
	var filePaths []string
	var files []*os.File

	fileNames := []string{
		LocalDockerComposeFileName,
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
