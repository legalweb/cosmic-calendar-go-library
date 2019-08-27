package commands

import (
	"fmt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/getopt"
)

type GetOAuthURLs struct {
	getopt.Command
	traits.Configurable
}

func NewGetOAuthURLs() (*GetOAuthURLs, error) {
	c := new(GetOAuthURLs)

	_, err := c.Command.Build("getoauthurls", c.Handle, nil)

	return c, err
}

func (c *GetOAuthURLs) Handle(opt *getopt.GetOpt) error {
	config, err := c.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	cs := calendar.NewCalendarService(config, false)
	r := cs.GetOAuthURLs()

	if len(r) > 0 {
		for k, v := range r {
			fmt.Println(k + " " + v)
		}
	} else {
		fmt.Println("No URLs retrieved.")
	}

	return nil
}
