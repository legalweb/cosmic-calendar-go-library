package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"google.golang.org/api/tasks/v1"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
)

type GetTasks struct {
	getopt.Command
	traits.Configurable
}

func NewTasks() (*GetTasks, error) {
	g := new(GetTasks)

	o1, _ := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")

	_, err := g.Command.Build("gettasks", g.Handle, o1)

	return g, err
}

func (g *GetTasks) Handle(opt *getopt.GetOpt) error {
	config, err := g.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	var r []*tasks.Task

	r, err = cs.GetTasks()

	if err != nil {
		return err
	}

	if r != nil {
		fmt.Println("Tasks")
		for _, e := range r {
			fmt.Println(*e)
		}
	} else {
		fmt.Println("No tasks retrieved.")
	}

	return nil
}
