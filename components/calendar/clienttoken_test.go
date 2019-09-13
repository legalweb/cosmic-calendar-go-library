package calendar

import (
	"encoding/json"
	"testing"
	"time"
)

func TestClientToken(t *testing.T) {
	x := NewClientToken()

	tokenStr := "Token String"
	vendorStr := "Vendor Name"

	jsonStr := "{\"expires\":\"2019-09-12T09:16:00Z\",\"token\":\"" + tokenStr + "\",\"vendor\":\"" + vendorStr + "\"}"
	err := json.Unmarshal([]byte(jsonStr), x)
	if err != nil {
		t.Error(err)
	}

	if x.Token != tokenStr {
		t.Errorf("got %q wanted %q", x.Token, tokenStr)
	}

	if x.Vendor != vendorStr {
		t.Errorf("got %q wanted %q", x.Vendor, vendorStr)
	}

	tp, err := time.Parse(time.RFC3339, "2019-09-12T09:16:00Z")
	if err != nil {
		t.Error(err)
	}

	if x.Expires != tp {
		t.Errorf("Time not converted correctly, got %q wanted %q", x.Expires, tp)
	}

	jsonStr = "{\"expires\":\"1234\",\"token\":\"Token String\",\"vendor\":\"Vendor Name\"}"
	err = json.Unmarshal([]byte(jsonStr), x)
	if err != nil {
		t.Error(err)
	}

	if x.Expires != time.Unix(1234, 0) {
		t.Errorf("Time not converted correctly, got %q wanted %q", x.Expires, time.Unix(1234, 0))
	}

	jsonStr = "{\"expires\":1234,\"token\":\"Token String\",\"vendor\":\"Vendor Name\"}"
	err = json.Unmarshal([]byte(jsonStr), x)
	if err != nil {
		t.Error(err)
	}

	if x.Expires != time.Unix(1234, 0) {
		t.Errorf("Time not converted correctly, got %q wanted %q", x.Expires, time.Unix(1234, 0))
	}

	jsonStr = "{\"expires\":null,\"token\":\"Token String\",\"vendor\":\"Vendor Name\"}"
	err = json.Unmarshal([]byte(jsonStr), x)
	if err == nil {
		t.Error("Error expected for value of expires")
	}

	jsonStr = "{\"expires\":1234,\"token\":1234,\"vendor\":\"Vendor Name\"}"
	err = json.Unmarshal([]byte(jsonStr), x)
	if err == nil {
		t.Error("Error expected for value of token")
	}

	jsonStr = "{\"expires\":1234,\"token\":\"Token String\",\"vendor\":1234}"
	err = json.Unmarshal([]byte(jsonStr), x)
	if err == nil {
		t.Error("Error expected for value of Vendor")
	}
}
