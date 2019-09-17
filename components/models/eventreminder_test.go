package models

import (
	"encoding/json"
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

func TestCanMarshalEventReminder(t *testing.T) {
	r := NewEventReminder("GET", 5)
	_, err := json.Marshal(r)

	if err != nil {
		t.Error(err)
	}
}
