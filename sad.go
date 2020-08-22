package sad

import "net"

type options struct {
	Server     net.IP
	Username   string
	RootDir    string
	PrivateKey interface{}
	Channel    string
	Path       string
	EnvVars    []string
	Debug      bool
}

// Hello says hello to the world
func Hello() string {
	return "Hello, world."
}
