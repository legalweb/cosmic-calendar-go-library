package calendar

import (
	"encoding/json"
	"time"
)

type ClientToken struct {
	Expires time.Time
	Token string
	Vendor string
}

func NewClientToken() *ClientToken {
	c := new(ClientToken)

	return c
}

func (c *ClientToken) UnmarshalJSON(data []byte) error {
	var pieces map[string]interface{}
	err := json.Unmarshal(data, &pieces)
	if err != nil {
		return err
	}

	for k, v := range pieces {
		switch k {
		case "expires": fallthrough
		case "Expires":
			c.Expires = time.Unix(int64(v.(float64)), 0)
		case "token": fallthrough
		case "Token":
			c.Token = v.(string)
		case "vendor": fallthrough
		case "Vendor":
			c.Vendor = v.(string)
		}
	}

	return nil
}