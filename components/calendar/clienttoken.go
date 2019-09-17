package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	_ = json.Unmarshal(data, &pieces)

	for k, v := range pieces {
		switch k {
		case "expires": fallthrough
		case "Expires":
			set := false
			switch v.(type) {
			case string:
				t, err := time.Parse(time.RFC3339, v.(string))
				if err != nil {
					i, err := strconv.ParseInt(v.(string), 10, 64)
					if err != nil {
						return err
					}
					t = time.Unix(i, 0)
				}
				c.Expires = t
				set = true
			case float64:
				c.Expires = time.Unix(int64(v.(float64)), 0)
				set = true
			}
			if !set {
				return errors.New(fmt.Sprintf("Unhandled type for Expires: %T", v))
			}
		case "token": fallthrough
		case "Token":
			set := false
			switch v.(type) {
			case string:
				c.Token = v.(string)
				set = true
			}
			if !set {
				return errors.New(fmt.Sprintf("Unhandled type for Token: %T", v))
			}
		case "vendor": fallthrough
		case "Vendor":
			set := false
			switch v.(type) {
			case string:
				c.Vendor = v.(string)
				set = true
			}
			if !set {
				return errors.New(fmt.Sprintf("Unhandled type for Vendor: %T", v))
			}
		}
	}

	return nil
}