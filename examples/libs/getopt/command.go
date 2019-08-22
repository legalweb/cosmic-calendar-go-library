package getopt

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	WithOptions
	WithOperands

	name string
	shortDescription string
	longDescription string

	handler handlerFunc
}

func NewCommand(name string, handler handlerFunc, options []Option) (*Command, error) {
	c := new(Command)
	_, err := c.SetName(name)

	if err != nil {
		return nil, err
	}

	c.handler = handler

	if len(options) > 0 {
		c.AddOptions(options)
	}

	return c, nil
}

func (c *Command) SetName(name string) (*Command, error) {
	if len(name) == 0 || name[0] == '-' || strings.Contains(name, " ") {
		return nil, errors.New(fmt.Sprintf("Command name has to be an alphanumeric string not starting with dash, found '%s'", name))
	}
	c.name = name
	return c, nil
}

func (c *Command) SetHandler(handler handlerFunc) *Command {
	c.handler = handler
	return c
}

func (c *Command) SetDescription(longDescription string) *Command {
	c.longDescription = longDescription
	if c.shortDescription == "" {
		c.shortDescription = longDescription
	}
	return c
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) GetHandler() *handlerFunc {
	return &c.handler
}

func (c *Command) GetDescription() string {
	return c.longDescription
}

func (c *Command) Description() string {
	return c.longDescription
}

func (c *Command) GetShortDescription() string {
	return c.shortDescription
}

func (c *Command) ShortDescription() string {
	return c.shortDescription
}