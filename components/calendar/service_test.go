package calendar

import (
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
	"reflect"
	"testing"
	"time"
)

var requester *MockCalendarRequester

func TestNewCalendarService(t *testing.T) {
	config := testConfig()

	x, err := Default()
	if x != nil || err == nil {
		t.Error("Calendar should be in unconfigured state")
	}

	requester = new(MockCalendarRequester)
	x = NewCalendarService(config, "1", true, requester)
	if reflect.TypeOf(x) != reflect.TypeOf(&CalendarService{}) {
		t.Errorf("Unexpected return type got %T expected %T", x, &CalendarService{})
	}

	x, err = Default()
	if x == nil || err != nil {
		t.Error("Calendar should be in configured state")
	}

	x, err = Default("Unconfigured")
	if x != nil || err == nil {
		t.Error("Calendar should be in unconfigured state")
	}

	x, err = Default("Test Name")
	if x == nil || err != nil {
		t.Error("Calendar should be in configured state")
	}
}

func TestGetClientToken(t *testing.T) {
	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	_, err = x.GetClientToken()
	if err != nil {
		t.Error(err)
	}

	requester.setJson("{\"ErrorMessage\":\"\",\"Response\":{\"NoToken\":\"\"},\"ResponseCode\": 200}", true)
	_, err = x.GetClientToken()
	if err == nil {
		t.Error("Request should return token not found error")
	}
}

func TestGetOAuthURLs(t *testing.T) {
	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	u, err := x.GetOAuthURLs()
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(u) != reflect.TypeOf(map[string]string{}) {
		t.Error("map of strings expected in return")
	}

	_, e := u["test"]
	if !e {
		t.Error("test url not found in return")
	}

	requester.setJson("{\"ErrorMessage\": \"\",\"Response\":{\"NoURLS\":\"\",\"Message\": \"OAuth URLs\"},\"ResponseCode\": 200}", true)
	_, err = x.GetOAuthURLs()
	if err == nil {
		t.Error("Request should return URLs not found error")
	}
}

func TestGetCalendlyLink(t *testing.T) {
	testUrl := "https://calendly.com/test-link"
	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	c, err := x.GetCalendlyLink()
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(c) != reflect.TypeOf(testUrl) {
		t.Error("url string expected in return")
	}

	if c != testUrl {
		t.Errorf("got %q wanted %q", c, testUrl)
	}

	requester.setJson("{\"ErrorMessage\":\"\",\"Response\":{\"NoUrl\":\"\",\"Message\":\"Calendly Link\"},\"ResponseCode\":200}", true)
	_, err = x.GetCalendlyLink()
	if err == nil {
		t.Error("Request should return url not found error")
	}
}

func TestGetEvents(t *testing.T) {
	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	c, err := x.GetEvents(5)
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(c) != reflect.TypeOf([]*calendar.Event{}) {
		t.Errorf("got %T wanted %T", c, []*calendar.Event{})
	}

	if len(c) != 1 {
		t.Error("Expected to receive 1 calendar event")
	}

	requester.setJson("{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Message\":\"Event List\"},\"ResponseCode\":200}", true)
	_, err = x.GetEvents()
	if err == nil {
		t.Error("Request should fail with no events in response")
	}

	requester.setJson("{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Events\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"p32gbbutksb6u80g\\\"\",\"kind\":\"calendar#events\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-13T12:11:15.556Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}", true)
	_, err = x.GetEvents()
	if err == nil {
		t.Error("Request should fail with no event items in response")
	}
}

func TestGetTasks(t *testing.T) {
	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	c, err := x.GetTasks()
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(c) != reflect.TypeOf([]*tasks.Task{}) {
		t.Errorf("got %T wanted %T", c, []*tasks.Task{})
	}

	if len(c) != 1 {
		t.Error("Expected to receive 1 calendar task")
	}

	requester.setJson("{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Message\":\"Event List\"},\"ResponseCode\":200}", true)
	_, err = x.GetTasks()
	if err == nil {
		t.Error("Request should fail with no tasks in response")
	}

	requester.setJson("{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Tasks\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"p32gbbutksb6u80g\\\"\",\"kind\":\"calendar#events\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-13T12:11:15.556Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}", true)
	_, err = x.GetTasks()
	if err == nil {
		t.Error("Request should fail with no task items in response")
	}
}

func TestSetCalendlyLink(t *testing.T) {
	testUrl := "https://calendly.com/test-link"
	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	c, err := x.SetCalendlyLink(testUrl)
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(c) != reflect.TypeOf(testUrl) {
		t.Error("url string expected in return")
	}

	if c != testUrl {
		t.Errorf("got %q wanted %q", c, testUrl)
	}
}

func TestAddEvent(t *testing.T) {
	summaryStr := "Test Event"
	descriptionStr := "Event Description"
	startStr := "2019-09-13T15:39:00Z"
	startTime, _ := time.Parse(time.RFC3339, startStr)
	endStr := "2019-09-13T16:39:00Z"
	endTime, _ := time.Parse(time.RFC3339, endStr)

	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	c, err := x.AddEvent(summaryStr, descriptionStr, startTime, endTime)
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(c) != reflect.TypeOf(&calendar.Event{}) {
		t.Errorf("got %T wanted %T", c, &calendar.Event{})
	}

	if c.Summary != summaryStr {
		t.Errorf("got %q wanted %q", c.Summary, summaryStr)
	}

	if c.Description != descriptionStr {
		t.Errorf("got %q wanted %q", c.Description, descriptionStr)
	}

	if c.Start.DateTime != startStr {
		t.Errorf("got %q wanted %q", c.Start.DateTime, startStr)
	}

	if c.End.DateTime != endStr {
		t.Errorf("got %q wanted %q", c.End.DateTime, endStr)
	}
}

func TestAddTask(t *testing.T) {
	titleStr := "Test Task"
	dueStr := "2019-09-14T00:00:00.000Z"
	dueTime, _ := time.Parse(time.RFC3339Nano, dueStr)

	requester.useJson(false)
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	c, err := x.AddTask(titleStr, dueTime)
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(c) != reflect.TypeOf(&tasks.Task{}) {
		t.Errorf("got %T wanted %T", c, &tasks.Task{})
	}

	if c.Title != titleStr {
		t.Errorf("got %q wanted %q", c.Title, titleStr)
	}

	if c.Due != dueStr {
		t.Errorf("got %q wanted %q", c.Due, dueStr)
	}
}

func testConfig() CalendarServiceConfig {
	c := CalendarServiceConfig{
		Name:      "Test Name",
		Client:    "Test Client",
		Secret:    "Test Secret",
		EndPoint:  "Test EndPoint",
		VerifySSL: false,
		Debug:     false,
	}

	return c
}
