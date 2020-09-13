package sad

import (
	"errors"
	"fmt"
	"os"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

// SendFiles sends files to a server using the provided SCP client.
func SendFiles(client *scp.Client, opts *Options, files []*os.File) error {
	err := client.Connect()

	if err != nil {
		return err
	}

	defer client.Close()

	for _, file := range files {
		defer file.Close()

		basename, err := file.Stat()

		if err != nil {
			errorMessage := fmt.Sprintf("Error stating file %s: %s", file.Name(), err)
			return errors.New(errorMessage)
		}

		remotePath := fmt.Sprintf("%s/%s/%s", opts.RootDir, GetFullName(opts), basename)
		permissions := "0655"
		err = client.CopyFile(file, remotePath, permissions)

		if err != nil {
			errorMessage := fmt.Sprintf("Error copying file %s to remote server: %s", file.Name(), err)
			return errors.New(errorMessage)
		}
	}

	return nil
}

// GetSSHClientConfig generates an SSH client config based on the provided options.
func GetSSHClientConfig(opts *Options) (*ssh.ClientConfig, error) {
	authMethod, err := opts.PrivateKey.ToSSHAuthMethod()

	if err != nil {
		return nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User:            opts.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{authMethod},
	}

	return clientConfig, nil
}

// GetSCPClient generates an SCP client based on the given options and SSH client config.
func GetSCPClient(opts *Options, clientConfig *ssh.ClientConfig) (*scp.Client, error) {
	port := 22
	host := fmt.Sprintf("%s:%d", opts.Server, port)

	duration, _ := time.ParseDuration("5s")

	scpClient := scp.NewClientWithTimeout(host, clientConfig, duration)

	return &scpClient, nil
}

// GetFullName gets the full name of the app given the provided options.
// The name is based on the app name and the channel.
func GetFullName(opts *Options) string {
	return fmt.Sprintf("%s-%s", opts.Name, opts.Channel)
}
