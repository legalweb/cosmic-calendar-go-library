package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestNewTaskRequest(t *testing.T) {
	title := "Task Title"
	due, _ := time.Parse(time.RFC3339, "2019-09-11T16:10:00Z")
	dueStr, _ := json.Marshal(due)
	jsonOutput := "{\"title\":\"" + title + "\",\"due\":" + string(dueStr) + "}"

	x := NewTaskRequest(title, due)

	if reflect.TypeOf(x) != reflect.TypeOf(&TaskRequest{}) {
		t.Errorf("got %s wanted %s", reflect.TypeOf(x).String(), reflect.TypeOf(&TaskRequest{}))
	}

	if x.Title != title {
		t.Errorf("got %q wanted %q", x.Title, title)
	}
	if x.Due != due {
		t.Errorf("got %q wanted %q", x.Due, due)
	}

	jStr, err := json.Marshal(x)

	if err != nil {
		t.Error(err)
	}

	if string(jStr) != jsonOutput {
		t.Errorf("got %q wanted %q", string(jStr), jsonOutput)
	}
}

func TestCanMarshalTaskRequest(t *testing.T) {
	title := "Task Title"
	due, _ := time.Parse(time.RFC3339, "2019-09-11T16:10:00Z")

	x := NewTaskRequest(title, due)

	_, err := json.Marshal(x)

	if err != nil {
		t.Error(err)
	}
}
