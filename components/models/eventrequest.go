package models

import (
	"encoding/json"
	"time"
)

type EventRequest struct {
	Summary string `json:"summary"`
	Description string `json:"description"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	Reminders []*EventReminder `json:"reminders"`
}

func NewEventRequest(summary string, start time.Time, end time.Time, reminders ...*EventReminder) *EventRequest {
	e := new(EventRequest)
	e.Summary = summary
	e.Start = start
	e.End = end
	if end.IsZero() {
		e.End = e.Start.Add(time.Minute * time.Duration(15))
	}
	e.Reminders = reminders

	return e
}

func (e *EventRequest) MarshalJSON() ([]byte, error) {
	summary, _ := json.Marshal(e.Summary)
	output := "{\"summary\":" + string(summary)

	if len(e.Description) > 0 {
		description, _ := json.Marshal(e.Description)
		output += ",\"description\":" + string(description)
	}

	start, _ := json.Marshal(e.Start)
	output += ",\"start\":" + string(start)

	if !e.End.IsZero() {
		end, _ := json.Marshal(e.End)
		output += ",\"end\":" + string(end)
	}

	if e.Reminders != nil && len(e.Reminders) > 0 {
		reminders, _ := json.Marshal(e.Reminders)
		output += ",\"reminders\":" + string(reminders)
	}

	return []byte(output + "}"), nil
}