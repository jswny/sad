package sad

import (
	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

func getSCPClient(opts *Options) (*scp.Client, error) {
	clientConfig, err := getClientConfig(opts)

	if err != nil {
		return nil, err
	}

	scpClient := scp.NewClient("example.com:22", clientConfig)

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
