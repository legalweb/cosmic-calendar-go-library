package models

import (
	"reflect"
	"testing"
)

func TestNewSetCalendlyLinkRequest(t *testing.T) {
	url := "https://test.com/"
	x := NewSetCalendlyLinkRequest(url)
	if reflect.TypeOf(x) != reflect.TypeOf(&SetCalendlyLinkRequest{}) {
		t.Errorf("got %s wanted %s", reflect.TypeOf(x).String(), reflect.TypeOf(&SetCalendlyLinkRequest{}))
	}
	if x.Url != url {
		t.Errorf("got %q wanted %q", x.Url, url)
	}
}
