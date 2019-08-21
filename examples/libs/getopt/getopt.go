package getopt

import (
	"os"
	"strconv"
)

const NO_ARGUMENT = ":noArg"
const REQUIRED_ARGUMENT = ":requiredArg"
const OPTIONAL_ARGUMENT = ":optionalArg"
const MULTIPLE_ARGUMENT = ":multipleArg"

const SETTING_SCRIPT_NAME = "scriptName"
const SETTING_DEFAULT_MODE = "defaultMode"
const SETTING_STRICT_OPTIONS = "strictOptions"
const SETTING_STRICT_OPERANCDS = "strictOperands"

var defaultTranslator *Translator

type Command interface {

}

type GetOpt struct {
	WithOperands
	WithOptions

	help Help
	settings map[string]bool
	operandsCount int
	commands map[string]*Command
	command *Command
	additionalOperands []string
	additionalOptions []Option
	translator Translator
}

func (g *GetOpt) Process(args ...string) error {
	if len(args) == 0 {
		args = os.Args
	}

}


func (g *GetOpt) getCommand(name string) *Command {
	if len(name) > 0 {
		v, isset := g.commands[name]
		if isset {
			return v
		}
		return nil
	}

	return g.command;
}

func (g *GetOpt) getCommands() map[string]*Command {
	return g.commands;
}

func (g *GetOpt) GetOperands() []string {
	operandsValues := make([]string, 0)

	for _, operand := range g.WithOperands.GetOperands() {
		value := operand.GetValue()

		if value == nil {
			continue
		}

		operandsValues = append(operandsValues, value...)
	}

	return append(operandsValues, g.additionalOperands...)
}

func (g *GetOpt) GetOperand(index string) string {
	operand := g.WithOperands.GetOperand(index)
	if operand != nil {
		return operand.GetValue()
	} else if (isInt(index)) {
		v, _ := strconv.Atoi(index)
		i := v - len(g.operands)
		ov, isset := g.additionalOperands[i]
		if i >= 0 && isset {
			return ov
		}
		return nil
	}
}

func Translate(key string) string {
	return GetTranslator().Translate(key)
}

func GetTranslator() *Translator {
	if defaultTranslator == nil {
		return new(Translator)
	}

	return defaultTranslator
}

func isInt(s string) bool {
	_, err := strconv.Atoi(s)

	if err == nil {
		return true
	}

	return false
}