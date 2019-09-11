package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
)

type GetCalendlyLink struct {
	getopt.Command
	traits.Configurable
}

func NewGetCalendlyLink() (*GetCalendlyLink, error) {
	c := new(GetCalendlyLink)

	o1, _ := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")

	_, err := c.Command.Build("getcalendlylink", c.Handle, o1)

	return c, err
}

func (c *GetCalendlyLink) Handle(opt *getopt.GetOpt) error {
	config, err := c.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	r, err := cs.GetCalendlyLink()

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
