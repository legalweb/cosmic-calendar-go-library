package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
	"time"
)

type AddTask struct {
	getopt.Command
	traits.Configurable
}

func NewAddTask() (*AddTask, error) {
	a := new(AddTask)

	o1, _ := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")
	o2, _ := getopt.NewOption('t', "title", getopt.REQUIRED_ARGUMENT, "Specify title of event")
	o3, _ := getopt.NewOption('d', "due", getopt.REQUIRED_ARGUMENT, "Specify due date time of task")

	_, err := a.Command.Build("addtask", a.Handle, o1, o2, o3)

	return a, err
}

func (a *AddTask) Handle(opt *getopt.GetOpt) error {
	config, err := a.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	title := opt.GetOptionString("title")
	due, err := time.Parse("2006-01-02T15:04:05Z", opt.GetOptionString("due"))

	if err != nil {
		return err
	}

	r, err := cs.AddTask(title, due)

	if err != nil {
		return err
	}

	if r != nil {
		fmt.Println("Task Added")
		fmt.Println(r)
	} else {
		fmt.Println("Task not added")
	}

	return nil
}
