package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/getopt"
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

	o1, err := getopt.NewOption('\x00', "version", getopt.NO_ARGUMENT)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
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

	_, err = opt.AddOptions([]*getopt.Option{
		o1, o2, o3, o4, o5, o6,
	})

	c1, err := getopt.NewCommand("test-setup", func(o *getopt.GetOpt) {
		fmt.Println("When you see this message the setup works." + "\n")
	}, nil)

	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	c1.SetDescription("Check if setup works")

	_, err = opt.AddCommand(c1)

	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	err = opt.Process()

	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	optValue := opt.GetOptionValue("version")
	fmt.Println(optValue)
	if optValue != nil && len(optValue) > 0 {
		fmt.Println(fmt.Sprintf("%s: %s", NAME, VERSION))
		os.Exit(0)
	}

	command := opt.GetCommand("")
	if command != nil || opt.GetOption("help") != nil {
		fmt.Println(opt.GetHelpText(nil))
		os.Exit(0)
	}

	handlerFunc := command.GetHandler()
	handlerFunc(opt)

	args := os.Args

	fmt.Println(args)
}