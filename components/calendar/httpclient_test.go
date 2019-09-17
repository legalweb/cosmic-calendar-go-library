package calendar

import (
	"net/http"
	"testing"
)

func TestMockHttpClient_Do(t *testing.T) {
	code := 200
	failcode := 501
	body := "Body content"

	client := new(mockHttpClient)
	client.SetResponse(code, body)

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Error(err)
	}

	var res interface{}

	res, err = client.Do(req)

	_, isResponse := res.(*http.Response)
	if !isResponse {
		t.Error("Do should return *http.Response")
	}
	if err != nil {
		t.Error(err)
	}

	client.SetResponse(failcode, body)

	res, err = client.Do(req)
	if err == nil {
		t.Error("Do should return mocked error for code 501")
	}
}

func TestMockHttpClient_SetResponse(t *testing.T) {
	code := 200
	body := "Body content"

	client := new(mockHttpClient)
	client.SetResponse(code, body)

	if client.code != code {
		t.Errorf("got %q wanted %q", client.code, code)
	}
	if client.body != body {
		t.Errorf("got %q wanted %q", client.body, body)
	}
}