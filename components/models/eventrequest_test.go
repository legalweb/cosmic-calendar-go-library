package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewEventRequest(t *testing.T) {
	summary := "Summary"
	description := "Description"
	start, _ := time.Parse(time.RFC3339, "2019-09-11T15:57:34Z")
	end, _ := time.Parse(time.RFC3339, "2019-10-11T15:57:34Z")
	zero := time.Time{}
	duration := time.Minute * time.Duration(15)
	jsonOutput := "{\"summary\":\"" + summary + "\",\"description\":\"" + description + "\",\"start\":\"2019-09-11T15:57:34Z\",\"end\":\"2019-10-11T15:57:34Z\",\"reminders\":[{\"method\":\"popup\",\"minutes\":15}]}"

	reminder := NewEventReminder("popup", 15)
	x := NewEventRequest(summary, description, start, end, reminder)

	if reflect.TypeOf(x) != reflect.TypeOf(&EventRequest{}) {
		t.Errorf("got %s wanted %s", reflect.TypeOf(x).String(), reflect.TypeOf(&EventRequest{}))
	}
	if x.Summary != summary {
		t.Errorf("got %q wanted %q", x.Summary, summary)
	}
	if x.Start != start {
		t.Errorf("got %q wanted %q", x.Start, start)
	}
	if x.End != end {
		t.Errorf("got %q wanted %q", x.End, end)
	}

	jStr, err := json.Marshal(x)

	if err != nil {
		t.Error(err)
	}

	if string(jStr) != jsonOutput {
		t.Errorf("got %q wanted %q", string(jStr), jsonOutput)
	}

	x = NewEventRequest(summary, description, start, zero)

	if reflect.TypeOf(x) != reflect.TypeOf(&EventRequest{}) {
		t.Errorf("got %s wanted %s", reflect.TypeOf(x).String(), reflect.TypeOf(&EventRequest{}))
	}
	if x.Summary != summary {
		t.Errorf("got %q wanted %q", x.Summary, summary)
	}
	if x.Start != start {
		t.Errorf("got %q wanted %q", x.Start, start)
	}
	if x.End != start.Add(duration) {
		t.Errorf("got %q wanted %q", x.End, end)
	}
}

func TestCanMarshalEventRequest(t *testing.T) {
	summary := "Summary"
	description := "Description"
	start, _ := time.Parse(time.RFC3339, "2019-09-11T15:57:34Z")
	end, _ := time.Parse(time.RFC3339, "2019-10-11T15:57:34Z")

	reminder := NewEventReminder("popup", 15)
	x := NewEventRequest(summary, description, start, end, reminder)

	_, err := json.Marshal(x)

	if err != nil {
		t.Error(err)
	}
}
