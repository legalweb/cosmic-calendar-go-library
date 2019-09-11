package commands

import (
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
)

type GetOAuthURLs struct {
	getopt.Command
	traits.Configurable
}

func NewGetOAuthURLs() (*GetOAuthURLs, error) {
	c := new(GetOAuthURLs)

	_, err := c.Command.Build("getoauthurls", c.Handle)

	return c, err
}

func (c *GetOAuthURLs) Handle(opt *getopt.GetOpt) error {
	config, err := c.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	cs := calendar.NewCalendarService(config, false)
	r, err := cs.GetOAuthURLs()

	if err != nil {
		return err
	}

	if len(r) > 0 {
		fmt.Println("URLS")
		for k, v := range r {
			fmt.Println(k + ": " + v)
		}
	} else {
		fmt.Println("No URLs retrieved.")
	}

	return nil
}
