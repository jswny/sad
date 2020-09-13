package sad

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

func sendFiles(client *scp.Client, opts *Options, files []*os.File) error {
	err := client.Connect()

	if err != nil {
		return err
	}

	defer client.Close()

	for _, file := range files {
		defer file.Close()

		remotePath := opts.RootDir
		permissions := "0655"
		err = client.CopyFile(file, remotePath, permissions)

		if err != nil {
			errorMessage := fmt.Sprintf("Error copying file %s to remote server: %s", file.Name(), err)
			return errors.New(errorMessage)
		}
	}

	return nil
}

func getSCPClient(opts *Options, clientConfig *ssh.ClientConfig) (*scp.Client, error) {
	port := 22
	host := fmt.Sprintf("%s:%s", opts.Server, strconv.Itoa(port))
	scpClient := scp.NewClient(host, clientConfig)

	return &scpClient, nil
}

func getClientConfig(opts *Options) (*ssh.ClientConfig, error) {
	authMethod, err := opts.PrivateKey.ToSSHAuthMethod()

	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User: opts.Username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
	}

	return clientConfig, nil
}
