package commands

import (
	"errors"
	"fmt"
	"github.com/daforester/getopt-golang/getopt"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
	"lwebco.de/cosmic-calendar-go-library/components/models"
	"lwebco.de/cosmic-calendar-go-library/examples/libs/commands/traits"
	"strconv"
	"time"
)

type AddEvent struct {
	getopt.Command
	traits.Configurable
}

func NewAddEvent() (*AddEvent, error) {
	a := new(AddEvent)

	o1, _ := getopt.NewOption('u', "user", getopt.REQUIRED_ARGUMENT, "Specify user to act on behalf of")
	o2, _ := getopt.NewOption('t', "title", getopt.REQUIRED_ARGUMENT, "Specify title of event")
	o3, _ := getopt.NewOption('\x00', "start", getopt.REQUIRED_ARGUMENT, "Specify start date time of event")
	o4, _ := getopt.NewOption('\x00', "end", getopt.OPTIONAL_ARGUMENT, "Specify end date time of event")
	o5, _ := getopt.NewOption('\x00', "email-reminder", getopt.OPTIONAL_ARGUMENT, "Minutes before event to send email reminder")
	o6, _ := getopt.NewOption('\x00', "popup-reminder", getopt.OPTIONAL_ARGUMENT, "Minutes before event to popup reminder")

	_, err := a.Command.Build("addevent", a.Handle, o1, o2, o3, o4, o5, o6)

	return a, err
}

func (a *AddEvent) Handle(opt *getopt.GetOpt) error {
	config, err := a.Configurable.GetCalendarConfig(opt)

	if err != nil {
		return err
	}

	user := opt.GetOptionString("user")

	if user == "" {
		return errors.New("User not configured for request")
	}

	cs := calendar.NewCalendarService(config, false, user)

	summary := opt.GetOptionString("title")
	start, err := time.Parse("2006-01-02T15:04:05Z", opt.GetOptionString("start"))
	var end time.Time

	if err != nil {
		return err
	}

	if opt.GetOptionString("end") != "" {
		end, err = time.Parse("2006-01-02T15:04:05Z", opt.GetOptionString("end"))

		if err != nil {
			return err
		}
	}

	var reminders []*models.EventReminder

	if mins := opt.GetOptionString("email-reminder"); mins != "" {
		if m, err := strconv.Atoi(mins); err == nil {
			reminders = append(reminders, models.NewEventReminder("email", m))
		}
	}

	if mins := opt.GetOptionString("popup-reminder"); mins != "" {
		if m, err := strconv.Atoi(mins); err == nil {
			reminders = append(reminders, models.NewEventReminder("popup", m))
		}
	}

	r, err := cs.AddEvent(summary, start, end, reminders...)

	if err != nil {
		return err
	}

	if r != nil {
		fmt.Println("Event Added")
		fmt.Println(r)
	} else {
		fmt.Println("Event not added")
	}

	return nil
}
