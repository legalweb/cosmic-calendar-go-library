package traits

import (
	"errors"
	"github.com/daforester/getopt-golang/getopt"
	"github.com/pelletier/go-toml"
	"lwebco.de/cosmic-calendar-go-library/components/calendar"
)

type Configurable struct {
}

func (c *Configurable) GetCalendarConfig(opt *getopt.GetOpt) (calendar.CalendarServiceConfig, error) {
	config, err := c.getConfigOptions(opt)

	if err != nil {
		return config, err
	}

	config.Name = "Example"

	return config, err
}

func (c *Configurable) getConfigOptions(opt *getopt.GetOpt) (calendar.CalendarServiceConfig, error) {
	calConfig := calendar.CalendarServiceConfig{}

	config := ""
	configValue := opt.GetOptionValue("config")

	if configValue != nil && len(configValue) > 0 {
		config = configValue[0]
	}

	client := ""
	secret := ""
	endpoint := ""
	verifyssl := false
	debug := false

	if config == "" {
		client = opt.GetOptionString("client")
		secret = opt.GetOptionString("secret")
		endpoint = opt.GetOptionString("endpoint")

		if client == "" || secret == "" || endpoint == "" {
			return calConfig, errors.New("client, secret & endpoint must be set if not providing config file")
		}
	} else {
		t, err := toml.LoadFile(config)
		if err != nil {
			return calConfig, err
		}

		if !t.Has("CLIENT") || !t.Has("SECRET") || !t.Has("ENDPOINT") {
			return calConfig, errors.New("client, secret & endpoint must be set in config file")
		}

		client = t.Get("CLIENT").(string)
		secret = t.Get("SECRET").(string)
		endpoint = t.Get("ENDPOINT").(string)
		verifyssl = t.Get("VERIFYSSL").(bool)
		debug = t.Get("DEBUG_REQUEST").(bool)
	}

	calConfig = calendar.CalendarServiceConfig{
		Client:    client,
		Secret:    secret,
		EndPoint:  endpoint,
		VerifySSL: verifyssl,
		Debug:     debug,
	}

	return calConfig, nil
}
