package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)


func main() {
	optVersion := getopt.BoolLong("version", 'v', "Show version information and quit")
	optHelp := getopt.BoolLong("help",'?', "Show this help and quit")
	optConfig := getopt.StringLong("config", 'f', "Specify configuration file to use")
	optClient := getopt.StringLong("client", 'c', "Provide client login if not config file provided")
	optSecret := getopt.StringLong("secret", 's', "Provide secret if no config file provided")
	optEndpoint := getopt.StringLong("endpoint", 'e', "Provide endpoint if no config file provided")

	getopt.Parse()

	getopt.CommandLine.Parse(args)

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	args := getopt.Args()

	fmt.Println(args)
}