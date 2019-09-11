package models

type EventReminder struct {
	Method string `json:"method"`
	Minutes int `json:"minutes"`
}

func NewEventReminder(method string, minutes int) *EventReminder {
	e := new(EventReminder)
	e.Method = method
	e.Minutes = minutes

	return e
}