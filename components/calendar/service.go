package calendar

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
	"io/ioutil"
	"lwebco.de/cosmic-calendar-go-library/components/models"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strconv"
	"time"
)

type CalendarService struct {
	config CalendarServiceConfig
	user   string
}

var (
	defaultCalendarService *CalendarService
	instances    map[string]*CalendarService
)

func NewCalendarService(config CalendarServiceConfig, opt ...interface{}) *CalendarService {
	m := new(CalendarService)
	m.SetConfig(config)

	isDefault := false
	user := ""

	for _, option := range opt {
		switch option.(type) {
		case bool:
			isDefault = option.(bool)
		case string:
			user = option.(string)
		}
	}

	if defaultCalendarService == nil || (isDefault == true) {
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
		return defaultCalendarService, errors.New("MySQL Database not configured")
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
	url := "/token/";

	r, err := s.request(url)

	if err != nil {
		return nil, err
	}

	if r["Token"] == nil {
		return nil, errors.New("Token not found in JSON response")
	}

	jStr, err := json.Marshal(r["Token"])

	if err != nil {
		return nil, err
	}

	clientToken := NewClientToken()

	err = json.Unmarshal([]byte(jStr), clientToken)

	if err != nil {
		return nil, err
	}

	return clientToken, nil
}

func (s *CalendarService) GetCalendlyLink() (string, error) {
	url := "/calendly/link";

	r, err := s.request(url)

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

	data, err := json.Marshal(setRequest)

	if err != nil {
		return "", err
	}

	url = "/calendly/link";

	r, err := s.request(url, string(data))

	if err != nil {
		return "", err
	}

	if r["Url"] == nil {
		return "", errors.New("Calendly link not set")
	}

	return r["Url"].(string), nil
}

func (s *CalendarService) AddEvent(summary string, start time.Time, end time.Time, reminders ...*models.EventReminder) (*calendar.Event, error) {
	eventRequest := models.NewEventRequest(summary, start, end, reminders...)

	data, err := json.Marshal(eventRequest)

	if err != nil {
		return nil, err
	}

	url := "/calendar/events";

	r, err := s.request(url, string(data))

	if err != nil {
		return nil, err
	}

	if r["Event"] == nil {
		return nil, errors.New("Event not created")
	}

	jStr, err := json.Marshal(r["Event"])

	if err != nil {
		return nil, err
	}

	event := new(calendar.Event)

	err = json.Unmarshal([]byte(jStr), event)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *CalendarService) AddTask(title string, due time.Time) (*tasks.Task, error) {
	taskRequest := models.NewTaskRequest(title, due)

	data, err := json.Marshal(taskRequest)

	if err != nil {
		return nil, err
	}

	url := "/calendar/tasks";

	r, err := s.request(url, string(data))

	if err != nil {
		return nil, err
	}

	if r["Task"] == nil {
		return nil, errors.New("Task not created")
	}

	jStr, err := json.Marshal(r["Task"])

	if err != nil {
		return nil, err
	}

	task := new(tasks.Task)

	err = json.Unmarshal([]byte(jStr), task)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *CalendarService) GetEvents(days ...int) ([]*calendar.Event, error) {
	noDays := 0

	url := "/calendar/events";

	if len(days) > 0 && days[0] > 0 {
		noDays = days[0]
	}

	if noDays > 0 {
		url += "?days" + strconv.Itoa(noDays)
	}

	r, err := s.request(url)

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

	jStr, err := json.Marshal(events["items"])

	if err != nil {
		return nil, err
	}

	var calEvents []*calendar.Event

	err = json.Unmarshal([]byte(jStr), &calEvents)

	if err != nil {
		return nil, err
	}

	return calEvents, nil
}

func (s *CalendarService) GetTasks() ([]*tasks.Task, error) {
	url := "/calendar/tasks";

	r, err := s.request(url)

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

	jStr, err := json.Marshal(events["items"])

	if err != nil {
		return nil, err
	}

	var calTasks []*tasks.Task

	err = json.Unmarshal([]byte(jStr), &calTasks)

	if err != nil {
		return nil, err
	}

	return calTasks, nil
}

func (s *CalendarService) GetOAuthURLs() (map[string]string, error) {
	r, err := s.request("/login/oauth/urls")

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

func (s *CalendarService) decodeResponse(j string) (map[string]interface{}, error) {

	type response struct {
		ErrorMessage string
		Response map[string]interface{}
		ResponseCode int
	}

	x := new(response)
	err := json.Unmarshal([]byte(j), x)

	if err != nil {
		return nil, err
	}

	if x.ResponseCode != 200 {
		return nil, errors.New(fmt.Sprintf("API Request failed: %d", x.ResponseCode))
	}

	return x.Response, err
}

func (s *CalendarService) request(url string, json ...string) (map[string]interface{}, error) {
	var req *http.Request
	var err error
	jsonString := ""

	if len(json) > 0 {
		for _, str := range json {
			jsonString += str
		}
	}

	buf := bytes.NewBuffer([]byte(jsonString))

	if buf.Len() > 0 {
		req, err = http.NewRequest("POST", s.config.EndPoint+url, buf)

		if err != nil {
			return nil, err
		}

		req.Header["Content-Type"] = []string{"application/json"}
	} else {
		req, err = http.NewRequest("GET", s.config.EndPoint+url, nil)

		if err != nil {
			return nil, err
		}
	}

	req.SetBasicAuth(s.config.Client, s.config.Secret)

	if len(s.user) > 0 {
		req.Header["X-Auth-User"] = []string{s.user}
	}

	if !s.config.VerifySSL {
		tr := http.DefaultTransport.(*http.Transport)

		if tr.TLSClientConfig == nil {
			tlsConf := new(tls.Config)
			tlsConf.InsecureSkipVerify = true
			tr.TLSClientConfig = tlsConf
		} else {
			tr.TLSClientConfig.InsecureSkipVerify = true
		}
	}

	if s.config.Debug {
		debug, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(debug))
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	output, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
		case 503:
			return nil, errors.New("API Service Unavailable")
		case 500:
			if s.config.Debug {
				fmt.Println("API Returned Error:")
				fmt.Println(string(output))
			}
			return nil, errors.New("API Returned Error")
		case 403:
			return nil, errors.New("Access Forbidden")
		case 400:
			return nil, errors.New("Bad API request")
		case 200:
			return s.decodeResponse(string(output))
	}

	return nil, errors.New("Unhandled API response: " + strconv.Itoa(res.StatusCode))
}