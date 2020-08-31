package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/jswny/sad"
)

func main() {
	fmt.Println("Sad CLI!")

	_, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseFlags() (sad.Options, error) {
	server := flag.String("server", "", "Server to deploy to")
	username := flag.String("username", "", "User to login to on the server")
	rootDir := flag.String("root-dir", "", "Root directory to deploy to on the server")
	privateKey := flag.String("private-key", "", "Base64 encoded SSH private key to login to the user on the server")
	channel := flag.String("channel", "", "Deployment channel")
	path := flag.String("path", "", "Path to the app to be deployed relative to the current directory")
	envVars := flag.String("env-vars", "", "Local environment variables to be injected into the app deployment")
	debug := flag.Bool("debug", false, "Debug mode")

	flag.Parse()

	opts := sad.Options{}
	debugString := strconv.FormatBool(*debug)
	err := opts.FromStrings(*server, *username, *rootDir, *privateKey, *channel, *path, *envVars, debugString)

	return opts, err
}
