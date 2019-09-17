package calendar

import (
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	config := testConfig()
	config.Debug = true

	x := NewCalendarService(config)

	client = &mockHttpClient{}
	requester = NewHTTPCalendarRequester(client)

	x = NewCalendarService(config, "1", true, requester)

	client.SetResponse(503, "")
	_, err := requester.Request(x, "/")
	if err == nil {
		t.Error("HTTP Request should return error for 503")
	}

	client.SetResponse(200, "{\"ErrorMessage\": \"\",\"Response\":{\"URLS\":{\"test\": \"https://example.com?oAuth\"},\"Message\": \"OAuth URLs\"},\"ResponseCode\": 200}")
	_, err = x.GetOAuthURLs()
	if err != nil {
		t.Error(err)
	}

	x.SetConfig(config)

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

	_, err = requester.Request(x, "invalidUrl:/x.net")

	if err == nil {
		t.Error("HTTP Request should return error for invalid GET URL")
	}

	_, err = requester.Request(x, "invalidUrl:/x.net", "{\"data\":\"foobar\"}")

	if err == nil {
		t.Error("HTTP Request should return error for invalid POST URL")
	}

	client.SetResponse(407, "")
	_, err = requester.Request(x, "/")

	if err == nil {
		t.Error("HTTP Request should return error for unknown errors")
	}

	client.SetResponse(501, "")
	_, err = requester.Request(x, "/")

	if err == nil {
		t.Error("HTTP Request should return error if request execution fails - mocked by 501 code")
	}

	var mockReadAll readAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, errors.New("Mocked error reading body")
	}

	var mockDumpRequest dumpRequestFunc = func(req *http.Request, body bool) ([]byte, error) {
		return nil, errors.New("Mocked error dumping request")
	}

	NewHTTPCalendarRequester(client, mockReadAll, mockDumpRequest)

	requester.setReadAllFunc(mockReadAll)

	client.SetResponse(200, "{\"ErrorMessage\": \"\",\"Response\":{\"URLS\":{\"test\": \"https://example.com?oAuth\"},\"Message\": \"OAuth URLs\"},\"ResponseCode\": 200}")
	_, err = x.GetOAuthURLs()
	if err == nil {
		t.Error("Expected error from readAll function")
	}

	requester.setDumpRequestFunc(mockDumpRequest)

	_, err = x.GetOAuthURLs()
	if err == nil {
		t.Error("Expected error from dumpRequest function")
	}

	defaultCalendarService = nil
}

func TestDecodeResponse(t *testing.T) {
	r := NewHTTPCalendarRequester()
	_, err := r.decodeResponse("")

	if err == nil {
		t.Error("Error expected decoding response")
	}

	_, err = r.decodeResponse("{\"ResponseCode\":501}")

	if err == nil {
		t.Error("API error expected")
	}
}
