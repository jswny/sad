package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jswny/sad"
	"golang.org/x/crypto/ssh"
)

var deploymentCommand string = "docker-compose up -d"

func main() {
	commandLineOpts, environmentOpts, configOpts := loadOptions()

	opts := checkOptions(commandLineOpts, environmentOpts, configOpts)

	clientConfig := configureSSHClient(opts)

	sshClient := openSSHConnection(clientConfig, opts)
	defer sshClient.Close()

	remotePath := getRemotePath(opts)

	createDeploymentDir(sshClient, remotePath)

	deployFiles(sshClient, opts)

	startApp(sshClient, remotePath, deploymentCommand)
}

// GetAllOptionSources gets options from each different source.
func GetAllOptionSources(program string, args []string, configFileName string) (commandLineOpts *sad.Options, environmentOpts *sad.Options, configOpts *sad.Options, commandLineOutput string, err error) {
	commandLineOpts, output, err := ParseFlags(program, args)
	if err != nil {
		return nil, nil, nil, output, err
	}

	environmentOpts = &sad.Options{}
	err = environmentOpts.FromEnv()
	if err != nil {
		return nil, nil, nil, "", err
	}

	configOpts = &sad.Options{}
	if configFileName != "" {
		err = configOpts.FromJSON(configFileName)
		if err != nil {
			return nil, nil, nil, "", err
		}
	}

	return commandLineOpts, environmentOpts, configOpts, "", nil
}

// MergeOptionsHierarchy merges options from different sources together.
// Options are merged in order starting from the options of least precedence to greatest precedence.
// Thus, the options with greatest precedence will contain the merged options.
// The sources in order of precedence are: command line, environment variables, config file.
func MergeOptionsHierarchy(commandLineOptions *sad.Options, environmentOptions *sad.Options, configOptions *sad.Options) {
	environmentOptions.Merge(configOptions)
	commandLineOptions.Merge(environmentOptions)
}

// ParseFlags parses command line flags into options.
// Flag parsing is always returned as output.
// If help or usage is requested, flag.ErrHelp is returned.
func ParseFlags(program string, args []string) (opts *sad.Options, output string, err error) {
	flags := flag.NewFlagSet(program, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	registry := flags.String("registry", "", "Docker image registry")
	image := flags.String("image", "", "Docker image to deploy")
	digest := flags.String("digest", "", "Docker image digest to deploy")
	server := flags.String("server", "", "Server to deploy to")
	username := flags.String("username", "", "User to login to on the server")
	rootDir := flags.String("root-dir", "", "Root directory to deploy to on the server")
	privateKey := flags.String("private-key", "", "Base64 encoded SSH private key to login to the user on the server")
	channel := flags.String("channel", "", "Deployment channel")
	envVars := flags.String("env-vars", "", "Local environment variables to be injected into the app deployment")
	debug := flags.Bool("debug", false, "Debug mode")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	opts = &sad.Options{}
	debugString := strconv.FormatBool(*debug)
	err = opts.FromStrings(*registry, *image, *digest, *server, *username, *rootDir, *privateKey, *channel, *envVars, debugString)

	if err != nil {
		return nil, buf.String(), err
	}

	return opts, buf.String(), nil
}

func loadOptions() (commandLineOpts *sad.Options, environmentOpts *sad.Options, configOpts *sad.Options) {
	fmt.Print("Loading config... ")

	configFilePath, err := sad.FindFilePathRecursive(".", sad.ConfigFileName)

	if err != nil {
		if err.Error() == sad.FindFilePathRecursiveFileNotFoundErrorMessage {
			fmt.Print("Could not find a config file, skipping... ")
		} else {
			fmt.Println("Error finding config file:", err)
			os.Exit(1)
		}
	} else {
		fmt.Print("Found config file: ", configFilePath, "... ")
	}

	commandLineOpts, environmentOpts, configOpts, commandLineOutput, err := GetAllOptionSources(os.Args[0], os.Args[1:], configFilePath)
	if err != nil {
		if commandLineOutput != "" {
			fmt.Println(commandLineOutput)
		}
		if err == flag.ErrHelp {
			os.Exit(2)
		}

		fmt.Println("Error retrieving options:", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
	return commandLineOpts, environmentOpts, configOpts
}

func checkOptions(commandLineOpts *sad.Options, environmentOpts *sad.Options, configOpts *sad.Options) *sad.Options {
	fmt.Print("Verifying config... ")

	MergeOptionsHierarchy(commandLineOpts, environmentOpts, configOpts)
	commandLineOpts.MergeDefaults()

	err := commandLineOpts.Verify()
	if err != nil {
		fmt.Println("Provided options were invalid:", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
	return commandLineOpts
}

func configureSSHClient(opts *sad.Options) *ssh.ClientConfig {
	fmt.Print("Configuring SSH client... ")

	clientConfig, err := sad.GetSSHClientConfig(opts)

	if err != nil {
		fmt.Println("Error getting SSH configuration from options:", err)
		os.Exit(1)
	}

	fmt.Println("Success!")

	return clientConfig
}

func openSSHConnection(clientConfig *ssh.ClientConfig, opts *sad.Options) *ssh.Client {
	fmt.Print("Opening SSH connection... ")

	address := opts.Server.String()
	port := "22"

	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(address, port), clientConfig)

	if err != nil {
		msg := fmt.Sprintf("failed to open SSH connection to address %s:%s: %s", address, port, err)
		fmt.Println(msg, err)
		os.Exit(1)
	}

	fmt.Println("Success!")

	return sshClient
}

func getRemotePath(opts *sad.Options) string {
	fmt.Print("Generating remote path... ")

	deploymentName, err := opts.GetDeploymentName()

	if err != nil {
		fmt.Println("Error getting full app name:", err)
		os.Exit(1)
	}

	remotePath := fmt.Sprintf("%s/%s", opts.RootDir, deploymentName)

	fmt.Println("Success!")

	return remotePath
}

func createDeploymentDir(sshClient *ssh.Client, remotePath string) {
	fmt.Print("Creating directory for deployment... ")

	cmd := fmt.Sprintf("mkdir -p %s", remotePath)
	output, err := sad.SSHRunCommand(sshClient, cmd)

	if err != nil {
		fmt.Println("Error creating directory for deployment:", err)
		maybePrettyPrintOutput(output)
		os.Exit(1)
	}

	maybePrettyPrintOutput(output)
	fmt.Println("Success!")
}

func deployFiles(sshClient *ssh.Client, opts *sad.Options) {
	fmt.Print("Sending files to server... ")

	readerMap, files, err := sad.GetEntitiesForDeployment(".", opts)

	for _, file := range files {
		defer file.Close()
	}

	if err != nil {
		fmt.Println("Error getting files for deployment:", err)
		os.Exit(1)
	}

	err = sad.SendFiles(sshClient, opts, readerMap)
	if err != nil {
		fmt.Println("Error sending files to server:", err)
		os.Exit(1)
	}

	fmt.Println("Success!")
}

func startApp(sshClient *ssh.Client, remotePath string, deploymentCommand string) {
	fmt.Print("Starting app on server... ")

	cmd := fmt.Sprintf("cd %s && %s", remotePath, deploymentCommand)

	output, err := sad.SSHRunCommand(sshClient, cmd)

	if err != nil {
		fmt.Println("Error starting app on server:", err)
		maybePrettyPrintOutput(output)
		os.Exit(1)
	}

	fmt.Println("Success!")

	maybePrettyPrintOutput(output)
}

func maybePrettyPrintOutput(output string) {
	lines := strings.Split(output, "\n")

	var prettyOutput string

	prefix := "|"
	for _, line := range lines {
		if line != "" {
			prettyOutput += prefix + " " + line + "\n"
		}
	}

	if prettyOutput != "" {
		fmt.Println(prettyOutput)
	}
}
