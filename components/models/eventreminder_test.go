package models

import (
	"reflect"
	"testing"
)

func TestNewEventReminder(t *testing.T) {
	x := NewEventReminder("GET", 5)
	if reflect.TypeOf(x) != reflect.TypeOf(&EventReminder{}) {
		t.Errorf("got %s wanted %s", reflect.TypeOf(x).String(), reflect.TypeOf(&EventReminder{}))
	}
	if x.Method != "GET" {
		t.Errorf("got %q wanted %q", x.Method, "GET")
	}
	if x.Minutes != 5 {
		t.Errorf("got %q wanted %q", x.Minutes, 5)
	}
}
