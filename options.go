package sad

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

// EnvVarPrefix represents the prefix that all environment variables should have to be read properly.
var EnvVarPrefix = "SAD_"

// Options for deployment.
type Options struct {
	Server     net.IP
	Username   string
	RootDir    string
	PrivateKey RSAPrivateKey
	Channel    string
	Path       string
	EnvVars    []string
	Debug      bool
}

// Merge merges the other options into the existing options
// When both fields are populated, the field from the existing options is kept.
func (o *Options) Merge(other *Options) {
	if o.Server == nil {
		o.Server = other.Server
	}

	if o.Username == "" {
		o.Username = other.Username
	}

	if o.RootDir == "" {
		o.RootDir = other.RootDir
	}

	if o.PrivateKey.PrivateKey == nil {
		o.PrivateKey = other.PrivateKey
	}

	if o.Channel == "" {
		o.Channel = other.Channel
	}

	if o.Path == "" {
		o.Path = other.Path
	}

	if len(o.EnvVars) == 0 {
		o.EnvVars = other.EnvVars
	}

	if !o.Debug {
		o.Debug = other.Debug
	}
}

// MergeDefaults merges default option values into the given options.
func (o *Options) MergeDefaults() {
	defaults := Options{
		Channel: "beta",
		Path:    ".",
		EnvVars: make([]string, 0),
		Debug:   false,
	}

	o.Merge(&defaults)
}

// FromStrings converts strings into options.
func (o *Options) FromStrings(server string, username string, rootDir string, privateKey string, channel string, path string, envVars string, debug string) error {
	if server != "" {
		o.Server = net.ParseIP(server)
	}

	o.Username = username
	o.RootDir = rootDir

	if privateKey != "" {
		rsaPrivateKey := RSAPrivateKey{}
		err := rsaPrivateKey.ParseBase64PEMString(privateKey)
		if err != nil {
			return err
		}
		o.PrivateKey = rsaPrivateKey
	}

	o.Channel = channel
	o.Path = path

	if envVars != "" {
		envVarsArr := strings.Split(envVars, ",")
		o.EnvVars = envVarsArr
	}

	if debug != "" {
		debugBool, err := strconv.ParseBool(debug)
		if err != nil {
			return err
		}

		o.Debug = debugBool
	}

	return nil
}

// GetJSON parses options from a JSON file.
func (o *Options) GetJSON(filename string) error {
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

// GetEnv parses options from environment variables.
// All variables should be prefixed and they should correspond to the available options with underscores separating words such as "PRIVATE_KEY".
// The private key should be a base64 encoded string.
// The environment variables should be a comma-separated string.
func (o *Options) GetEnv() error {
	prefix := EnvVarPrefix

	server := os.Getenv(prefix + "SERVER")
	username := os.Getenv(prefix + "USERNAME")
	rootDir := os.Getenv(prefix + "ROOT_DIR")
	privateKey := os.Getenv(prefix + "PRIVATE_KEY")
	channel := os.Getenv(prefix + "CHANNEL")
	path := os.Getenv(prefix + "PATH")
	envVars := os.Getenv(prefix + "ENV_VARS")
	debug := os.Getenv(prefix + "DEBUG")

	err := o.FromStrings(server, username, rootDir, privateKey, channel, path, envVars, debug)

	if err != nil {
		return err
	}

	return nil
}
