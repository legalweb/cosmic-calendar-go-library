package calendar

type MockCalendarRequester struct {
	returnError bool
	useProvidedJson bool
	json string
	HTTPCalendarRequester
}

func NewMockCalendarRequester() *MockCalendarRequester {
	return new(MockCalendarRequester)
}

func (m *MockCalendarRequester) Request(s *CalendarService, url string, json ...string) (map[string]interface{}, error) {
	jStr := ""

	switch (url) {
	case "/token": fallthrough
	case "/token/":
		jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Token\":{\"Expires\":1568316330,\"Token\":\"kK0=\",\"Vendor\":\"Testing\"}},\"ResponseCode\": 200}"
	case "/login/oauth/urls": fallthrough
	case "/login/oauth/urls/":
		jStr = "{\"ErrorMessage\": \"\",\"Response\":{\"URLS\":{\"test\": \"https://example.com?oAuth\"},\"Message\": \"OAuth URLs\"},\"ResponseCode\": 200}"
	case "/api/calendly/link": fallthrough
	case "/api/calendly/link/":
		jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Url\":\"https://calendly.com/test-link\",\"Message\":\"Calendly Link\"},\"ResponseCode\":200}"
	case "/api/calendar/events": fallthrough
	case "/api/calendar/events/":
		jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@example.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@example.com\",\"timeZone\":\"Europe/London\"},\"Events\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"etag\\\"\",\"kind\":\"calendar#events\",\"summary\":\"example@example.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-12T14:49:13.911Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}"
	case "/api/calendar/tasks": fallthrough
	case "/api/calendar/tasks/":
		jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"List\":{\"etag\":\"\\\"etag\\\"\",\"id\":\"id\",\"kind\":\"tasks#taskList\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/users/@me/lists/url\",\"title\":\"Task List\",\"updated\":\"2019-09-12T14:51:11.928Z\"},\"Tasks\":{\"etag\":\"\\\"etag\\\"\",\"kind\":\"tasks#tasks\"},\"Message\":\"Task List\"},\"ResponseCode\":200}"
	}

	if m.useProvidedJson {
		return m.decodeResponse(m.json)
	}

	return m.decodeResponse(jStr)
}

func (m *MockCalendarRequester) setJson(s string, use ...bool) {
	m.json = s
	m.useJson(true)

	if use != nil && len(use) > 0 {
		m.useJson(use[0])
	}
}

func (m *MockCalendarRequester) useJson(b bool) {
	m.useProvidedJson = b
}
