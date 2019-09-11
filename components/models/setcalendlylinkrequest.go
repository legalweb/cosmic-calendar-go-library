package models

type SetCalendlyLinkRequest struct {
	Url string `json:"url"`
}

func NewSetCalendlyLinkRequest(url string) *SetCalendlyLinkRequest {
	e := new(SetCalendlyLinkRequest)
	e.Url = url

	return e
}
