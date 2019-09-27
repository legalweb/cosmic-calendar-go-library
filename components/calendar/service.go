package calendar

import (
	"encoding/json"
	"errors"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
	"lwebco.de/cosmic-calendar-go-library/components/models"
	"reflect"
	"strconv"
	"time"
)

type CalendarService struct {
	config    CalendarServiceConfig
	user      string
	requester CalendarRequester
}

var (
	defaultCalendarService *CalendarService
	instances              map[string]*CalendarService
)

func NewCalendarService(config CalendarServiceConfig, opt ...interface{}) *CalendarService {
	m := new(CalendarService)
	m.SetConfig(config)

	isDefault := false
	user := ""

	for _, option := range opt {
		switch option := option.(type) {
		case bool:
			isDefault = option
		case string:
			user = option
		}
		_, isRequester := interface{}(option).(CalendarRequester)
		if isRequester {
			m.requester = option.(CalendarRequester)
		}
	}

	if m.requester == nil {
		m.requester = NewHTTPCalendarRequester()
	}

	if defaultCalendarService == nil || isDefault {
		defaultCalendarService = m
	}

	if len(config.Name) > 0 {
		if instances == nil {
			instances = make(map[string]*CalendarService)
		}

		instances[config.Name] = m
	}

	if len(user) > 0 {
		m.SetUser(user)
	}

	return m
}

func Default(name ...string) (*CalendarService, error) {
	if len(name) > 0 && len(name[0]) > 0 {
		m := instances[name[0]]
		if m == nil {
			return nil, errors.New("CalendarService " + name[0] + " not configured")
		}

		return m, nil
	}

	if defaultCalendarService == nil {
		return defaultCalendarService, errors.New("Calendar Service not configured")
	}

	return defaultCalendarService, nil
}

func (s *CalendarService) SetConfig(config CalendarServiceConfig) {
	s.config = config
}

func (s *CalendarService) SetUser(user string) {
	s.user = user
}

func (s *CalendarService) GetClientToken() (*ClientToken, error) {
	url := "/token"

	r, err := s.requester.Request(s, url)

	if err != nil {
		return nil, err
	}

	if r["Token"] == nil {
		return nil, errors.New("Token not found in JSON response")
	}

	clientToken := NewClientToken()
	err = s.remarshal(r["Token"], clientToken)

	if err != nil {
		return nil, err
	}

	return clientToken, nil
}

func (s *CalendarService) GetCalendlyLink() (string, error) {
	url := "/calendly/link"

	r, err := s.requester.Request(s, url)

	if err != nil {
		return "", err
	}

	if r["Url"] == nil {
		return "", errors.New("Calendly link not set")
	}

	return r["Url"].(string), nil
}

func (s *CalendarService) SetCalendlyLink(url string) (string, error) {
	setRequest := models.NewSetCalendlyLinkRequest(url)

	data, _ := json.Marshal(setRequest)

	url = "/calendly/link"

	r, err := s.requester.Request(s, url, string(data))

	if err != nil {
		return "", err
	}

	if r["Url"] == nil {
		return "", errors.New("Calendly link not set")
	}

	return r["Url"].(string), nil
}

func (s *CalendarService) AddEvent(summary string, description string, start time.Time, end time.Time, reminders ...*models.EventReminder) (*calendar.Event, error) {
	eventRequest := models.NewEventRequest(summary, description, start, end, reminders...)

	data, _ := json.Marshal(eventRequest)

	url := "/calendar/events"

	r, err := s.requester.Request(s, url, string(data))

	if err != nil {
		return nil, err
	}

	if r["Event"] == nil {
		return nil, errors.New("Event not created")
	}

	event := new(calendar.Event)
	err = s.remarshal(r["Event"], event)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *CalendarService) AddTask(title string, due time.Time) (*tasks.Task, error) {
	taskRequest := models.NewTaskRequest(title, due)

	data, _ := json.Marshal(taskRequest)

	url := "/calendar/tasks"

	r, err := s.requester.Request(s, url, string(data))

	if err != nil {
		return nil, err
	}

	if r["Task"] == nil {
		return nil, errors.New("Task not created")
	}

	task := new(tasks.Task)
	err = s.remarshal(r["Task"], task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *CalendarService) GetEvents(days ...int) ([]*calendar.Event, error) {
	noDays := 0

	url := "/calendar/events"

	if len(days) > 0 && days[0] > 0 {
		noDays = days[0]
	}

	if noDays > 0 {
		url += "?days=" + strconv.Itoa(noDays)
	}

	r, err := s.requester.Request(s, url)

	if err != nil {
		return nil, err
	}

	if r["Events"] == nil {
		return nil, errors.New("Events not found in JSON response")
	}

	events := r["Events"].(map[string]interface{})

	if events["items"] == nil {
		return nil, errors.New("Event items not found in JSON response")
	}

	var calEvents []*calendar.Event
	err = s.remarshal(events["items"], &calEvents)

	if err != nil {
		return nil, err
	}

	return calEvents, nil
}

func (s *CalendarService) GetTasks() ([]*tasks.Task, error) {
	url := "/calendar/tasks"

	r, err := s.requester.Request(s, url)

	if err != nil {
		return nil, err
	}

	if r["Tasks"] == nil {
		return nil, errors.New("Tasks not found in JSON response")
	}

	events := r["Tasks"].(map[string]interface{})

	if events["items"] == nil {
		return nil, errors.New("Task items not found in JSON response")
	}

	var calTasks []*tasks.Task
	err = s.remarshal(events["items"], &calTasks)

	if err != nil {
		return nil, err
	}

	return calTasks, nil
}

func (s *CalendarService) GetOAuthURLs() (map[string]string, error) {
	r, err := s.requester.Request(s, "/login/oauth/urls")

	if err != nil {
		return nil, err
	}

	urlsObj, exists := r["URLS"].(map[string]interface{})

	if !exists {
		return nil, errors.New("URLs object not found in JSON Response")
	}

	urls := make(map[string]string, len(urlsObj))

	for k, v := range urlsObj {
		if reflect.TypeOf(v).Kind() == reflect.String {
			urls[k] = v.(string)
		}
	}

	return urls, nil
}

func (s *CalendarService) remarshal(input interface{}, output interface{}) error {
	jStr, err := json.Marshal(input)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jStr, output)

	if err != nil {
		return err
	}

	return nil
}
