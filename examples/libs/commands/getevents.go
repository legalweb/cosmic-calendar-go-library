package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	calendar2 "google.golang.org/api/calendar/v3"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
	"strconv"
)

type GetEvents struct {
	getopt.Command
	traits.Configurable
}

func NewEvents() (*GetEvents, error) {
	g := new(GetEvents)

	o1, _ := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")
	o2, _ := getopt.NewOption('d', "days", getopt.OPTIONAL_ARGUMENT, "Specify number of days ahead to obtain events for")

	_, err := g.Command.Build("getevents", g.Handle, o1, o2)

	return g, err
}

func (g *GetEvents) Handle(opt *getopt.GetOpt) error {
	config, err := g.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	var r []*calendar2.Event

	days := opt.GetOptionString("days")
	if len(days) > 0 {
		i, err := strconv.Atoi(days)
		if err != nil {
			return err
		}
		r, err = cs.GetEvents(i)
	} else {
		r, err = cs.GetEvents()
	}

	if err != nil {
		return err
	}

	if r != nil {
		fmt.Println("Events")
		for _, e := range r {
			fmt.Println(*e)
		}
	} else {
		fmt.Println("No events retrieved.")
	}

	return nil
}
