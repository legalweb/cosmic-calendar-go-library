package models

import (
	"encoding/json"
	"time"
)

type TaskRequest struct {
	Title string    `json:"title"`
	Due   time.Time `json:"due"`
}

func NewTaskRequest(title string, due time.Time) *TaskRequest {
	e := new(TaskRequest)
	e.Title = title
	e.Due = due

	return e
}

func (e *TaskRequest) MarshalJSON() ([]byte, error) {
	title, _ := json.Marshal(e.Title)
	output := "{\"title\":" + string(title)

	due, _ := json.Marshal(e.Due)
	output += ",\"due\":" + string(due)

	return []byte(output + "}"), nil
}
