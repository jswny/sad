package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jswny/sad"
)

func main() {
	_, output, err := ParseFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		fmt.Println(output)
		os.Exit(2)
	} else if err != nil {
		fmt.Println("Error parsing command line arguments: ", err)
		fmt.Println(output)
		os.Exit(1)
	}

	log.Println("Starting deployment...")
}

// ParseFlags parses command line flags into options.
// Flag parsing is always returned as output.
// If help or usage is requested, flag.ErrHelp is returned.
func ParseFlags(program string, args []string) (opts *sad.Options, output string, err error) {
	flags := flag.NewFlagSet(program, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	server := flags.String("server", "", "Server to deploy to")
	username := flags.String("username", "", "User to login to on the server")
	rootDir := flags.String("root-dir", "", "Root directory to deploy to on the server")
	privateKey := flags.String("private-key", "", "Base64 encoded SSH private key to login to the user on the server")
	channel := flags.String("channel", "", "Deployment channel")
	path := flags.String("path", "", "Path to the app to be deployed relative to the current directory")
	envVars := flags.String("env-vars", "", "Local environment variables to be injected into the app deployment")
	debug := flags.Bool("debug", false, "Debug mode")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	opts = &sad.Options{}
	debugString := strconv.FormatBool(*debug)
	err = opts.FromStrings(*server, *username, *rootDir, *privateKey, *channel, *path, *envVars, debugString)

	if err != nil {
		return nil, buf.String(), err
	}

	return opts, buf.String(), nil
}
