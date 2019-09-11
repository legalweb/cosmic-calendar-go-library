package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
)

type SetCalendlyLink struct {
	getopt.Command
	traits.Configurable
}

func NewSetCalendlyLink() (*SetCalendlyLink, error) {
	c := new(SetCalendlyLink)

	o1, _ := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")
	o2, _ := getopt.NewOption('l', "url", getopt.REQUIRED_ARGUMENT, "Specify URL to calendly account")

	_, err := c.Command.Build("setcalendlylink", c.Handle, o1, o2)

	return c, err
}

func (c *SetCalendlyLink) Handle(opt *getopt.GetOpt) error {
	config, err := c.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	url := opt.GetOptionString("url")

	r, err := cs.SetCalendlyLink(url)

	if err != nil {
		return err
	}

	if r != "" {
		fmt.Println("Calendly Link")
		fmt.Println(r)
	} else {
		fmt.Println("No link set.")
	}

	return nil
}
