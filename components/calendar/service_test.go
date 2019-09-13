package calendar

import (
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
	"reflect"
	"testing"
	"time"
)

var client *mockHttpClient
var requester *HTTPCalendarRequester

func TestNewCalendarService(t *testing.T) {
	config := testConfig()

	x, err := Default()
	if x != nil || err == nil {
		t.Error("Calendar should be in unconfigured state")
	}

	client = &mockHttpClient{}
	requester = new(HTTPCalendarRequester)
	requester.SetClient(client)

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
	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Token\":{\"Expires\":1568316330,\"Token\":\"kK0=\",\"Vendor\":\"Testing\"}},\"ResponseCode\": 200}")
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	_, err = x.GetClientToken()
	if err != nil {
		t.Error(err)
	}
	
	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"NoToken\":\"\"},\"ResponseCode\": 200}")
	_, err = x.GetClientToken()
	if err == nil {
		t.Error("Request should return token not found error")
	}
}

func TestGetOAuthURLs(t *testing.T) {
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\": \"\",\"Response\":{\"URLS\":{\"test\": \"https://example.com?oAuth\"},\"Message\": \"OAuth URLs\"},\"ResponseCode\": 200}")
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
	
	client.SetResponse(200, "{\"ErrorMessage\": \"\",\"Response\":{\"NoURLS\":\"\",\"Message\": \"OAuth URLs\"},\"ResponseCode\": 200}")
	_, err = x.GetOAuthURLs()
	if err == nil {
		t.Error("Request should return URLs not found error")
	}
}

func TestGetCalendlyLink(t *testing.T) {
	testUrl := "https://calendly.com/test-link"
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Url\":\"https://calendly.com/test-link\",\"Message\":\"Calendly Link\"},\"ResponseCode\":200}")
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

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"NoUrl\":\"\",\"Message\":\"Calendly Link\"},\"ResponseCode\":200}")
	_, err = x.GetCalendlyLink()
	if err == nil {
		t.Error("Request should return url not found error")
	}
}

func TestGetEvents(t *testing.T) {
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Events\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"p32gbbutksb6u80g\\\"\",\"items\":[{\"created\":\"2019-09-13T12:11:15.000Z\",\"creator\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"end\":{\"dateTime\":\"2019-09-14T14:45:00+01:00\"},\"etag\":\"\\\"3136753351112000\\\"\",\"htmlLink\":\"https://www.google.com/calendar/event?eid=eid\",\"iCalUID\":\"9tlha4bu7ol6e5jpg82o4l39t0@google.com\",\"id\":\"9tlha4bu7ol6e5jpg82o4l39t0\",\"kind\":\"calendar#event\",\"organizer\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"reminders\":{\"useDefault\":true},\"start\":{\"dateTime\":\"2019-09-14T14:30:00+01:00\"},\"status\":\"confirmed\",\"summary\":\"Test Event\",\"updated\":\"2019-09-13T12:11:15.556Z\"}],\"kind\":\"calendar#events\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-13T12:11:15.556Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}")
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

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Message\":\"Event List\"},\"ResponseCode\":200}")
	_, err = x.GetEvents()
	if err == nil {
		t.Error("Request should fail with no events in response")
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Events\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"p32gbbutksb6u80g\\\"\",\"kind\":\"calendar#events\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-13T12:11:15.556Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}")
	_, err = x.GetEvents()
	if err == nil {
		t.Error("Request should fail with no event items in response")
	}
}

func TestGetTasks(t *testing.T) {
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"List\":{\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/V99u5YpzjWv20L6cnXYkJWaWzwc\\\"\",\"id\":\"MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA\",\"kind\":\"tasks#taskList\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/users/@me/lists/MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA\",\"title\":\"Cosmic\",\"updated\":\"2019-09-13T13:44:28.948Z\"},\"Tasks\":{\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/NzE5MjA2Mjk2\\\"\",\"items\":[{\"due\":\"2019-09-14T00:00:00.000Z\",\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/NzE5MjA2MjY1\\\"\",\"id\":\"QXFqSWJSR0ktUm85eTVMTQ\",\"kind\":\"tasks#task\",\"position\":\"00000000000000000000\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/lists/MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA/tasks/QXFqSWJSR0ktUm85eTVMTQ\",\"status\":\"needsAction\",\"title\":\"Test Task\",\"updated\":\"2019-09-13T13:44:28.000Z\"}],\"kind\":\"tasks#tasks\"},\"Message\":\"Task List\"},\"ResponseCode\":200}")
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

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Message\":\"Event List\"},\"ResponseCode\":200}")
	_, err = x.GetTasks()
	if err == nil {
		t.Error("Request should fail with no tasks in response")
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Tasks\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"p32gbbutksb6u80g\\\"\",\"kind\":\"calendar#events\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-13T12:11:15.556Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}")
	_, err = x.GetTasks()
	if err == nil {
		t.Error("Request should fail with no task items in response")
	}
}

func TestSetCalendlyLink(t *testing.T) {
	testUrl := "https://calendly.com/test-link"
	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Url\":\"https://calendly.com/test-link\",\"Message\":\"Calendly Link\"},\"ResponseCode\":200}")
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
	startStr := "2019-09-13T16:45:00+01:00"
	startTime, _ := time.Parse(time.RFC3339, startStr)
	endStr := "2019-09-13T16:55:00+01:00"
	endTime, _ := time.Parse(time.RFC3339, endStr)

	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Event\":{\"created\":\"2019-09-13T14:45:51.000Z\",\"creator\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"end\":{\"dateTime\":\"2019-09-13T16:55:00+01:00\"},\"etag\":\"\\\"3136771902110000\\\"\",\"htmlLink\":\"https://www.google.com/calendar/event?eid=b2dwOTJzdG12NnEzMHNsczFwMHJ0b2czc2sgYWFyb24ucGFya2VyQGIyYmZpbmFuY2UuY29t\",\"iCalUID\":\"ogp92stmv6q30sls1p0rtog3sk@google.com\",\"id\":\"ogp92stmv6q30sls1p0rtog3sk\",\"kind\":\"calendar#event\",\"organizer\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"reminders\":{\"useDefault\":true},\"start\":{\"dateTime\":\"2019-09-13T16:45:00+01:00\"},\"status\":\"confirmed\",\"summary\":\"Test Event\",\"description\":\"Event Description\",\"updated\":\"2019-09-13T14:45:51.055Z\"},\"Message\":\"Event Created\"},\"ResponseCode\":200}")
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

	x, err := Default()
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(200, "{\"ErrorMessage\":\"\",\"Response\":{\"Task\":{\"due\":\"2019-09-14T00:00:00.000Z\",\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/NzIzNjQ1NDI3\\\"\",\"id\":\"Rl94MlE0NkVaVGpERUlQag\",\"kind\":\"tasks#task\",\"position\":\"00000000000000000000\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/lists/MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA/tasks/Rl94MlE0NkVaVGpERUlQag\",\"status\":\"needsAction\",\"title\":\"Test Task\",\"updated\":\"2019-09-13T14:58:28.000Z\"},\"Message\":\"Task Created\"},\"ResponseCode\":200}")
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
