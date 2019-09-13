package calendar

import (
	"reflect"
	"testing"
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
