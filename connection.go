package sad

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

// SendFiles sends the specified reader interfaces as files to a server using the provided SCP client.
// The files are specified as a map of the name of the file to send to the server to a reader which can read the file.
// The full path name for the file on the remote server will be generatd as <root directory as specified by options>/<app name with channel>/<file name>.
func SendFiles(client *scp.Client, opts *Options, files map[string]io.Reader) error {
	err := client.Connect()

	if err != nil {
		return err
	}

	defer client.Close()

	for fileName, reader := range files {
		remotePath := fmt.Sprintf("%s/%s/%s", opts.RootDir, opts.GetFullAppName(), fileName)
		permissions := "0655"
		err = client.CopyFile(reader, remotePath, permissions)

		if err != nil {
			return fmt.Errorf("Error copying file %s to remote server: %s", fileName, err)
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

// SSHRunCommand Runs the specified command via SSH given the server to connect to and a client configuration.
// Returns the output of the command, or an error.
func SSHRunCommand(address string, port string, clientConfig *ssh.ClientConfig, cmd string) (string, error) {
	client, err := ssh.Dial("tcp", net.JoinHostPort(address, port), clientConfig)

	if err != nil {
		return "", fmt.Errorf("failed to dial SSH connection to address %s:%s: %s", address, port, err)
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	err = session.Run(cmd)

	if err != nil {
		return "", fmt.Errorf("failed to execute command \"%s\" via SSH connection to address %s:%s: %s", cmd, address, port, err)
	}

	return b.String(), nil
}
