package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/jswny/sad"
)

func main() {
	configFilePath, err := sad.FindFilePathRecursive(".", sad.ConfigFileName)

	if err != nil {
		if err.Error() == sad.FindFilePathRecursiveFileNotFoundErrorMessage {
			fmt.Println("Could not find a config file, skipping...")
		} else {
			fmt.Println("Error finding config file: ", err)
			os.Exit(1)
		}
	}

	commandLineOpts, environmentOpts, configOpts, commandLineOutput, err := GetAllOptionSources(os.Args[0], os.Args[1:], configFilePath)
	if err != nil {
		if commandLineOutput != "" {
			fmt.Println(commandLineOutput)
		}
		if err == flag.ErrHelp {
			os.Exit(2)
		}

		fmt.Println("Error retrieving options: ", err)
		os.Exit(1)
	}

	fmt.Println("Starting deployment...")

	MergeOptionsHierarchy(commandLineOpts, environmentOpts, configOpts)
	commandLineOpts.MergeDefaults()

	err = commandLineOpts.Verify()
	if err != nil {
		fmt.Println("Provided options were invalid!")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Opening SCP connection...")
	clientConfig, err := sad.GetSSHClientConfig(commandLineOpts)

	if err != nil {
		fmt.Println("Error getting SSH configuration from options: ", err)
		os.Exit(1)
	}

	scpClient, err := sad.GetSCPClient(commandLineOpts, clientConfig)

	if err != nil {
		fmt.Println("Error getting SCP client: ", err)
		os.Exit(1)
	}

	fmt.Println("Sending files to server...")

	files, err := sad.GetFilesForDeployment(".")

	if err != nil {
		fmt.Println("Eror getting files for deployment: ", err)
		os.Exit(1)
	}

	for _, file := range files {
		defer file.Close()
	}

	readerMap := sad.FilesToFileNameReaderMap(files)
	readerMap[sad.RemoteDockerComposeFileName] = readerMap[sad.LocalDockerComposeFileName]
	delete(readerMap, sad.LocalDockerComposeFileName)

	env := commandLineOpts.GetEnvValues()
	readerMap[sad.RemoteDotEnvFileName] = sad.GenerateDotEnvFile(env)

	err = sad.SendFiles(scpClient, commandLineOpts, readerMap)
	if err != nil {
		fmt.Println("Error sending files to server: ", err)
		os.Exit(1)
	}

	fmt.Println("Opening SSH connection...")
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
	err = configOpts.FromJSON(configFileName)
	if err != nil {
		return nil, nil, nil, "", err
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

	name := flags.String("name", "", "Name of the app to deploy")
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
	err = opts.FromStrings(*name, *server, *username, *rootDir, *privateKey, *channel, *envVars, debugString)

	if err != nil {
		return nil, buf.String(), err
	}

	return opts, buf.String(), nil
}
