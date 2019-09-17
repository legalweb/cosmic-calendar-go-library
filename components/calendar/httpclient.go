package calendar

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type mockHttpClient struct {
	http.Client
	body string
	code int
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	r := new(http.Response)
	r.Body = ioutil.NopCloser(strings.NewReader(m.body))
	r.StatusCode = m.code

	if m.code == 501 {
		return nil, errors.New("Mocked Response: Method Not Implemented")
	}

	return r, nil
}

func (m *mockHttpClient) SetResponse(code int, body string) {
	m.code = code
	m.body = body
}
