package main

import (
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"github.com/daforester/getopt-golang/getopt/errors"
	"github.com/sirupsen/logrus"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands"
	"os"
)

const NAME = "run"
const VERSION = "1.0-alpha"

func main() {
	opt, err := getopt.NewGetOpt("", nil)

	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	addOptions(opt)
	addCommands(opt)

	err = opt.Process()

	if err != nil {
		logrus.Errorln(err)
		if err.(errors.GetOptError).Type() == errors.ERROR_UNEXPECTED_ARGUMENT {
			fmt.Println("\n" + opt.GetHelpText(nil))
		}
		os.Exit(1)
	}

	optValue := opt.GetOptionValue("version")
	if optValue != nil && len(optValue) > 0 {
		fmt.Println(fmt.Sprintf("%s: %s", NAME, VERSION))
		os.Exit(0)
	}

	command := opt.GetCommand("")
	optHelp := opt.GetOptionValue("help")
	if command == nil || (optHelp != nil && len(optHelp) > 0) {
		fmt.Println(opt.GetHelpText(nil))
		os.Exit(0)
	}

	handlerFunc := command.GetHandler()
	if handlerFunc != nil {
		err := handlerFunc(opt)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Error("No handlerFunc found")
	}
}

func addCommands(opt *getopt.GetOpt) {
	c1, _ := getopt.NewCommand("test-setup", func(o *getopt.GetOpt) error {
		fmt.Println("When you see this message the setup works." + "\n")
		return nil
	})
	c1.SetDescription("Check if setup works")

	c2, _ := commands.NewGetOAuthURLs()
	c3, _ := commands.NewEvents()
	c4, _ := commands.NewAddEvent()
	c5, _ := commands.NewTasks()
	c6, _ := commands.NewAddTask()
	c7, _ := commands.NewGetClientToken()
	c8, _ := commands.NewGetCalendlyLink()
	c9, err := commands.NewSetCalendlyLink()

	if err != nil {
		logrus.Error(err)
	}

	_, _ = opt.AddCommands(c1, c2, c3, c4, c5, c6, c7, c8, c9)
}

func addOptions(opt *getopt.GetOpt) {
	o1, _ := getopt.NewOption('\x00', "version", getopt.NO_ARGUMENT)
	o1.SetDescription("Show version information and quit")
	o2, _ := getopt.NewOption('?', "help", getopt.NO_ARGUMENT)
	o2.SetDescription("Show this help and quit")
	o3, _ := getopt.NewOption('f', "config", getopt.OPTIONAL_ARGUMENT)
	o3.SetValidation(getopt.FileIsReadable, nil)
	o3.SetDescription("Specify configuration file to use")
	o4, _ := getopt.NewOption('c', "client", getopt.OPTIONAL_ARGUMENT)
	o4.SetDescription("Provide client login if no config file provided")
	o5, _ := getopt.NewOption('s', "secret", getopt.OPTIONAL_ARGUMENT)
	o5.SetDescription("Provide secret if no config file provided")
	o6, _ := getopt.NewOption('e', "endpoint", getopt.OPTIONAL_ARGUMENT)
	o6.SetDescription("Provide endpoint if no config file provided")

	_, _ = opt.AddOptions(o1, o2, o3, o4, o5, o6)
}