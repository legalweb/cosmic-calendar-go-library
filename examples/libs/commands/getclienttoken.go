package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
)

type GetClientToken struct {
	getopt.Command
	traits.Configurable
}

func NewGetClientToken() (*GetClientToken, error) {
	c := new(GetClientToken)

	o1, err := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")

	_, err = c.Command.Build("getclienttoken", c.Handle, o1)

	return c, err
}

func (c *GetClientToken) Handle(opt *getopt.GetOpt) error {
	config, err := c.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	r, err := cs.GetClientToken()

	if err != nil {
		return err
	}

	if r != nil {
		fmt.Println("Token")
		fmt.Println(r)
	} else {
		fmt.Println("No token retrieved.")
	}

	return nil
}
