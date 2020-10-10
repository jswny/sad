package sad

import (
	"errors"
	"os"
	"path/filepath"
)

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
