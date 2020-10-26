package sad

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// EnvVarPrefix represents the prefix that all environment variables should have to be read properly.
var EnvVarPrefix = "SAD_"

// Options for deployment.
type Options struct {
	Repository  string
	Server      net.IP
	Username    string
	RootDir     string
	PrivateKey  RSAPrivateKey
	Channel     string
	EnvVars     []string
	Debug       bool
	ImageDigest string
}

// Merge merges the other options into the existing options
// When both fields are populated, the field from the existing options is kept.
func (o *Options) Merge(other *Options) {
	if o.Repository == "" {
		o.Repository = other.Repository
	}

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

	if len(o.EnvVars) == 0 {
		o.EnvVars = other.EnvVars
	}

	if !o.Debug {
		o.Debug = other.Debug
	}

	if o.ImageDigest == "" {
		o.ImageDigest = other.ImageDigest
	}
}

// MergeDefaults merges default option values into the given options.
func (o *Options) MergeDefaults() {
	defaults := Options{
		Channel: "beta",
		Debug:   false,
	}

	o.Merge(&defaults)
}

// Verify verifies that the options are valid.
// Returns an error with information about why the options are invalid.
func (o *Options) Verify() error {
	errorMap := make(map[string]string)
	empty := "<empty>"

	if o.Repository == "" {
		errorMap["repository"] = fmt.Sprintf("is %s", empty)
	}

	if o.Server == nil {
		errorMap["server"] = "is nil"
	}

	if o.Username == "" {
		errorMap["username"] = fmt.Sprintf("is %s", empty)
	}

	if o.RootDir == "" {
		errorMap["root directory"] = fmt.Sprintf("is %s", empty)
	}

	if o.PrivateKey.PrivateKey == nil {
		errorMap["private key"] = "is nil"
	}

	if o.Channel == "" {
		errorMap["channel"] = fmt.Sprintf("is %s", empty)
	}

	if o.ImageDigest == "" {
		errorMap["image digest"] = fmt.Sprintf("is %s", empty)
	}

	if len(errorMap) != 0 {
		errorString := "invalid options! "

		for field, message := range errorMap {
			errorString += fmt.Sprintf("%s %s, ", field, message)
		}

		errorString = errorString[:len(errorString)-2]

		return fmt.Errorf(errorString)
	}

	return nil
}

// FromStrings converts strings into options.
func (o *Options) FromStrings(repository string, server string, username string, rootDir string, privateKey string, channel string, envVars string, debug string, imageDigest string) error {
	o.Repository = repository

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

	o.ImageDigest = imageDigest

	return nil
}

// FromJSON parses options from a JSON file.
func (o *Options) FromJSON(path string) error {
	file, err := ioutil.ReadFile(path)

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

// FromEnv parses options from environment variables.
// All variables should be prefixed and they should correspond to the available options with underscores separating words such as "PRIVATE_KEY".
// The private key should be a base64 encoded string.
// The environment variables should be a comma-separated string.
func (o *Options) FromEnv() error {
	prefix := EnvVarPrefix

	repository := os.Getenv(prefix + "REPOSITORY")
	server := os.Getenv(prefix + "SERVER")
	username := os.Getenv(prefix + "USERNAME")
	rootDir := os.Getenv(prefix + "ROOT_DIR")
	privateKey := os.Getenv(prefix + "PRIVATE_KEY")
	channel := os.Getenv(prefix + "CHANNEL")
	envVars := os.Getenv(prefix + "ENV_VARS")
	debug := os.Getenv(prefix + "DEBUG")
	imageDigest := os.Getenv(prefix + "IMAGE_DIGEST")

	err := o.FromStrings(repository, server, username, rootDir, privateKey, channel, envVars, debug, imageDigest)

	if err != nil {
		return err
	}

	return nil
}

// GetDeploymentName gets the full name of the deployment given the provided options.
// The name is based on the repository and the channel.
// All non-alphanumeric characters are replaced by dashes.
func (o *Options) GetDeploymentName() (string, error) {
	regStr := "[^a-zA-Z0-9]+"
	reg, err := regexp.Compile(regStr)

	if err != nil {
		return "", fmt.Errorf("error compiling regex %s: %s", regStr, err)
	}

	deploymentName := fmt.Sprintf("%s-%s", o.Repository, o.Channel)
	deploymentName = reg.ReplaceAllString(deploymentName, "-")

	return deploymentName, nil
}

// GetEnvValues gets the values of the environment variables specified in the EnvVars field.
// Returns a map of the variable names to values.
func (o *Options) GetEnvValues() map[string]string {
	m := make(map[string]string)

	for _, variableName := range o.EnvVars {
		value := os.Getenv(variableName)

		m[variableName] = value
	}

	return m
}
