package calendar

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (s *CalendarService) GetClientToken() {

}

func (s *CalendarService) GetCalendlyLink() {

}

func (s *CalendarService) SetCalendlyLink(url string) {

}

func (s *CalendarService) AddEvent(summary string, start time.Time, end time.Time) {

}

func (s *CalendarService) AddTask(title string, due time.Time) {

}

func (s *CalendarService) GetEvents(days ...int) (interface{}, error) {
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

	fmt.Println("Success getting events")
	fmt.Println(r)

	return r, nil
}

func (s *CalendarService) GetTasks() {

}

func (s *CalendarService) GetOAuthURLs() {

}

func (s *CalendarService) mustHaveUser() {

}

func (s *CalendarService) decodeResponse() {

}

func (s *CalendarService) request(url string, json ...string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.SetBasicAuth(s.config.Client, s.config.Secret)

	if len(s.user) > 0 {
		req.Header["X-Auth-User"] = []string{s.user}
	}

	if len(json) > 0 {
		req.Method = "POST"
		req.Header["Content-Type"] = []string{"application/json"}
		req.Header["Content-Length"] = []string{strconv.Itoa(len(json))}
		for _, d := range json {
			_, err = req.Body.Read([]byte(d))
			if err != nil {
				return "", err
			}
		}
	}

	if !s.config.VerifySSL {
		tr := http.DefaultTransport.(*http.Transport)
		tr.TLSClientConfig.InsecureSkipVerify = true
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	return string(body), nil
}