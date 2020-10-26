package sad

import (
	"bytes"
	"fmt"
	"io"

	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

// SendFiles sends the specified reader interfaces as files to a server using the provided SSH client.
// The files are specified as a map of the name of the file to send to the server to a reader which can read the file.
// The full path name for the file on the remote server will be generatd as <root directory as specified by options>/<app name with channel>/<file name>.
func SendFiles(sshClient *ssh.Client, opts *Options, files map[string]io.Reader) error {
	for fileName, reader := range files {
		fullName, err := opts.GetFullAppName()

		if err != nil {
			return fmt.Errorf("error getting full app name: %s", err)
		}

		remotePath := fmt.Sprintf("%s/%s/%s", opts.RootDir, fullName, fileName)

		permissions := "0644"
		err = copyFile(fileName, reader, remotePath, permissions, sshClient)

		if err != nil {
			return err
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

// SSHRunCommand Runs the specified command via SSH given the specified client.
// Returns the output of the command, or an error.
func SSHRunCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()

	if err != nil {
		return "", err
	}

	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)

	output := stdout.String() + stderr.String()

	if err != nil {
		return output, fmt.Errorf("failed to execute command \"%s\" via SSH client: %s", cmd, err)
	}

	return output, nil
}

func copyFile(fileName string, reader io.Reader, remotePath string, permissions string, sshClient *ssh.Client) error {
	client, err := scp.NewClientBySSH(sshClient)

	if err != nil {
		return fmt.Errorf("error creating new SSH session for SCP using existing SSH connection: %s", err)
	}

	defer client.Close()

	err = client.CopyFile(reader, remotePath, permissions)

	if err != nil {
		return fmt.Errorf("error copying file %s to remote server: %s", fileName, err)
	}

	return nil
}
