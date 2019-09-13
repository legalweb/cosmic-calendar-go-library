package calendar

import "testing"

func TestRequest(t *testing.T) {
	config := testConfig()

	x := NewCalendarService(config)

	client = &mockHttpClient{}
	requester = NewHTTPCalendarRequester(client)

	x = NewCalendarService(config, "1", true, requester)

	client.SetResponse(503, "")
	_, err := requester.Request(x, "/")
	if err == nil {
		t.Error("HTTP Request should return error for 503")
	}
	client.SetResponse(500, "")
	_, err = requester.Request(x, "/")
	if err == nil {
		t.Error("HTTP Request should return error for 500")
	}
	client.SetResponse(403, "")
	_, err = requester.Request(x, "/")

	if err == nil {
		t.Error("HTTP Request should return error for 403")
	}
	client.SetResponse(400, "")
	_, err = requester.Request(x, "/")

	if err == nil {
		t.Error("HTTP Request should return error for 400")
	}

	defaultCalendarService = nil
}