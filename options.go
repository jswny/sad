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

// OptionEnvVarPrefix represents the prefix that all environment variables representing options should have to be read properly.
var OptionEnvVarPrefix = "SAD_"

// DeploymentEnvVarPrefix represents the prefix that all dynamic environment variables which will be injected into the deployment should have to be read properly.
var DeploymentEnvVarPrefix = OptionEnvVarPrefix + "DEPLOY_"

// Options for deployment.
type Options struct {
	Registry   string
	Image      string
	Digest     string
	Server     net.IP
	Username   string
	RootDir    string
	PrivateKey RSAPrivateKey
	Channel    string
	EnvVars    []string
	Debug      bool
}

// Merge merges the other options into the existing options
// When both fields are populated, the field from the existing options is kept.
func (o *Options) Merge(other *Options) {
	if o.Registry == "" {
		o.Registry = other.Registry
	}

	if o.Image == "" {
		o.Image = other.Image
	}

	if o.Digest == "" {
		o.Digest = other.Digest
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
}

// MergeDefaults merges default option values into the given options.
func (o *Options) MergeDefaults() {
	defaults := Options{
		Channel: "beta",
		RootDir: "/",
		Debug:   false,
	}

	o.Merge(&defaults)
}

// Verify verifies that the options are valid.
// Returns an error with information about why the options are invalid.
func (o *Options) Verify() error {
	errorMap := make(map[string]string)
	empty := "<empty>"

	if o.Image == "" {
		errorMap["image"] = fmt.Sprintf("is %s", empty)
	}

	if o.Digest == "" {
		errorMap["digest"] = fmt.Sprintf("is %s", empty)
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
func (o *Options) FromStrings(registry string, image string, digest string, server string, username string, rootDir string, privateKey string, channel string, envVars string, debug string) error {
	o.Registry = registry

	o.Image = image

	o.Digest = digest

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
	prefix := OptionEnvVarPrefix

	registry := os.Getenv(prefix + "REGISTRY")
	image := os.Getenv(prefix + "IMAGE")
	digest := os.Getenv(prefix + "DIGEST")
	server := os.Getenv(prefix + "SERVER")
	username := os.Getenv(prefix + "USERNAME")
	rootDir := os.Getenv(prefix + "ROOT_DIR")
	privateKey := os.Getenv(prefix + "PRIVATE_KEY")
	channel := os.Getenv(prefix + "CHANNEL")
	envVars := os.Getenv(prefix + "ENV_VARS")
	debug := os.Getenv(prefix + "DEBUG")

	err := o.FromStrings(registry, image, digest, server, username, rootDir, privateKey, channel, envVars, debug)

	if err != nil {
		return err
	}

	return nil
}

// GetDeploymentName gets the full name of the deployment.
// The name is based on the image and the channel.
// All non-alphanumeric characters are replaced by dashes.
func (o *Options) GetDeploymentName() (string, error) {
	deploymentName := fmt.Sprintf("%s-%s", o.Image, o.Channel)
	deploymentName, err := replaceNonAlphanumeric(deploymentName, "-")

	if err != nil {
		return "", fmt.Errorf("error replacing non-alphanumeric characters in deployment name: %s", err)
	}

	return deploymentName, nil
}

// GetImageSpecifier gets the full image specifier for the deployment.
// The specifier is based on the image and the digest.
func (o *Options) GetImageSpecifier() string {
	separator := ""

	if o.Registry != "" {
		separator = "/"
	}

	specifier := fmt.Sprintf("%s%s%s@%s", o.Registry, separator, o.Image, o.Digest)

	return specifier
}

// GetDeploymentEnvValues gets the values of the environment variables specified in the EnvVars field to be injected into the deployment.
// Returns a map of the variable names to values, or an error if any of the variables are blank or unset.
func (o *Options) GetDeploymentEnvValues() (map[string]string, error) {
	m := make(map[string]string)

	for _, variableName := range o.EnvVars {
		variableNameWithPrefix := DeploymentEnvVarPrefix + variableName
		value := os.Getenv(variableNameWithPrefix)

		if value == "" {
			return nil, fmt.Errorf("environment variable %s is blank or unset", variableNameWithPrefix)
		}

		m[variableName] = value
	}

	return m, nil
}

func replaceNonAlphanumeric(input string, replaceWith string) (string, error) {
	regStr := "[^a-zA-Z0-9]+"
	reg, err := regexp.Compile(regStr)

	if err != nil {
		return "", fmt.Errorf("error compiling regex %s: %s", regStr, err)
	}

	return reg.ReplaceAllString(input, replaceWith), nil
}
