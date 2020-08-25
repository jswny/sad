package sad

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
)

// Options for deployment
type Options struct {
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

// Get parses options from a file
func (o *Options) Get(filename string) error {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, o)
}
