package calendar

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
)

type dumpRequestFunc func(*http.Request, bool) ([]uint8, error)
type readAllFunc func(io.Reader) ([]uint8, error)

type CalendarRequester interface {
	Request(*CalendarService, string, ...string) (map[string]interface{}, error)
	decodeResponse(string) (map[string]interface{}, error)
}

type HTTPCalendarRequester struct {
	client      httpClient
	dumpRequest dumpRequestFunc
	readAll     readAllFunc
}

func NewHTTPCalendarRequester(o ...interface{}) *HTTPCalendarRequester {
	r := new(HTTPCalendarRequester)

	r.client = http.DefaultClient
	r.dumpRequest = httputil.DumpRequest
	r.readAll = ioutil.ReadAll

	if len(o) > 0 {
		for _, obj := range o {

			switch obj := obj.(type) {
			case httpClient:
				r.setClient(obj)
			case dumpRequestFunc:
				r.setDumpRequestFunc(obj)
			case readAllFunc:
				r.setReadAllFunc(obj)
			}
		}
	}

	return r
}

func (r *HTTPCalendarRequester) setClient(c httpClient) {
	r.client = c
}

func (r *HTTPCalendarRequester) setDumpRequestFunc(dumpRequest dumpRequestFunc) {
	r.dumpRequest = dumpRequest
}

func (r *HTTPCalendarRequester) setReadAllFunc(realAll readAllFunc) {
	r.readAll = realAll
}

func (r *HTTPCalendarRequester) Request(s *CalendarService, url string, json ...string) (map[string]interface{}, error) {
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
		debug, err := r.dumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(debug))
	}

	res, err := r.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	output, err := r.readAll(res.Body)

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
		return r.decodeResponse(string(output))
	}

	return nil, errors.New("Unhandled API response: " + strconv.Itoa(res.StatusCode))
}

func (r *HTTPCalendarRequester) decodeResponse(j string) (map[string]interface{}, error) {

	type response struct {
		ErrorMessage string
		Response     map[string]interface{}
		ResponseCode int
	}

	x := new(response)
	err := json.Unmarshal([]byte(j), x)

	if err != nil {
		return nil, err
	}

	if x.ResponseCode != 200 {
		return nil, fmt.Errorf("API Request failed: %d", x.ResponseCode)
	}

	return x.Response, err
}
