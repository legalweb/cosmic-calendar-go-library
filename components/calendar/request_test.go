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
	case "/calendly/link": fallthrough
	case "/calendly/link/":
		jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Url\":\"https://calendly.com/test-link\",\"Message\":\"Calendly Link\"},\"ResponseCode\":200}"
	case "/calendar/events?days=5": fallthrough
	case "/calendar/events/?days=5": fallthrough
	case "/calendar/events": fallthrough
	case "/calendar/events/":
		if len(json) > 0 {
			jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Event\":{\"created\":\"2019-09-13T14:45:51.000Z\",\"creator\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"end\":{\"dateTime\":\"2019-09-13T16:55:00+01:00\"},\"etag\":\"\\\"3136771902110000\\\"\",\"htmlLink\":\"https://www.google.com/calendar/event?eid=b2dwOTJzdG12NnEzMHNsczFwMHJ0b2czc2sgYWFyb24ucGFya2VyQGIyYmZpbmFuY2UuY29t\",\"iCalUID\":\"ogp92stmv6q30sls1p0rtog3sk@google.com\",\"id\":\"ogp92stmv6q30sls1p0rtog3sk\",\"kind\":\"calendar#event\",\"organizer\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"reminders\":{\"useDefault\":true},\"start\":{\"dateTime\":\"2019-09-13T16:45:00+01:00\"},\"status\":\"confirmed\",\"summary\":\"Test Event\",\"updated\":\"2019-09-13T14:45:51.055Z\"},\"Message\":\"Event Created\"},\"ResponseCode\":200}"
		} else {
			jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Calendar\":{\"conferenceProperties\":{\"allowedConferenceSolutionTypes\":[\"hangoutsMeet\"]},\"etag\":\"\\\"etag\\\"\",\"id\":\"example@b2bfinance.com\",\"kind\":\"calendar#calendar\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\"},\"Events\":{\"accessRole\":\"owner\",\"defaultReminders\":[{\"method\":\"popup\",\"minutes\":10}],\"etag\":\"\\\"p32gbbutksb6u80g\\\"\",\"items\":[{\"created\":\"2019-09-13T12:11:15.000Z\",\"creator\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"end\":{\"dateTime\":\"2019-09-14T14:45:00+01:00\"},\"etag\":\"\\\"3136753351112000\\\"\",\"htmlLink\":\"https://www.google.com/calendar/event?eid=eid\",\"iCalUID\":\"9tlha4bu7ol6e5jpg82o4l39t0@google.com\",\"id\":\"9tlha4bu7ol6e5jpg82o4l39t0\",\"kind\":\"calendar#event\",\"organizer\":{\"email\":\"example@b2bfinance.com\",\"self\":true},\"reminders\":{\"useDefault\":true},\"start\":{\"dateTime\":\"2019-09-14T14:30:00+01:00\"},\"status\":\"confirmed\",\"summary\":\"Test Event\",\"updated\":\"2019-09-13T12:11:15.556Z\"}],\"kind\":\"calendar#events\",\"summary\":\"example@b2bfinance.com\",\"timeZone\":\"Europe/London\",\"updated\":\"2019-09-13T12:11:15.556Z\"},\"Message\":\"Event List\"},\"ResponseCode\":200}"
		}
	case "/calendar/tasks": fallthrough
	case "/calendar/tasks/":
		if len(json) > 0 {
			jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"Task\":{\"due\":\"2019-09-14T00:00:00.000Z\",\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/NzIzNjQ1NDI3\\\"\",\"id\":\"Rl94MlE0NkVaVGpERUlQag\",\"kind\":\"tasks#task\",\"position\":\"00000000000000000000\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/lists/MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA/tasks/Rl94MlE0NkVaVGpERUlQag\",\"status\":\"needsAction\",\"title\":\"Test Task\",\"updated\":\"2019-09-13T14:58:28.000Z\"},\"Message\":\"Task Created\"},\"ResponseCode\":200}"
		} else {
			jStr = "{\"ErrorMessage\":\"\",\"Response\":{\"List\":{\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/V99u5YpzjWv20L6cnXYkJWaWzwc\\\"\",\"id\":\"MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA\",\"kind\":\"tasks#taskList\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/users/@me/lists/MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA\",\"title\":\"Cosmic\",\"updated\":\"2019-09-13T13:44:28.948Z\"},\"Tasks\":{\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/NzE5MjA2Mjk2\\\"\",\"items\":[{\"due\":\"2019-09-14T00:00:00.000Z\",\"etag\":\"\\\"8MEupY6AVkDup3m0O6mGjMTpTY8/NzE5MjA2MjY1\\\"\",\"id\":\"QXFqSWJSR0ktUm85eTVMTQ\",\"kind\":\"tasks#task\",\"position\":\"00000000000000000000\",\"selfLink\":\"https://www.googleapis.com/tasks/v1/lists/MTEwMjE0NjA4NjE4ODc5OTcwMDg6MDEyMDQzNDE1ODY1OTA2NTgxMzU6MA/tasks/QXFqSWJSR0ktUm85eTVMTQ\",\"status\":\"needsAction\",\"title\":\"Test Task\",\"updated\":\"2019-09-13T13:44:28.000Z\"}],\"kind\":\"tasks#tasks\"},\"Message\":\"Task List\"},\"ResponseCode\":200}"
		}
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
