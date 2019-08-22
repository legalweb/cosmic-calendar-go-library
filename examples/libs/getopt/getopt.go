package getopt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const NO_ARGUMENT = ":noArg"
const REQUIRED_ARGUMENT = ":requiredArg"
const OPTIONAL_ARGUMENT = ":optionalArg"
const MULTIPLE_ARGUMENT = ":multipleArg"

var SETTING_COMMAND_NAME = "scriptName"
var SETTING_DEFAULT_MODE = "defaultMode"
var SETTING_STRICT_OPTIONS = "strictOptions"
var SETTING_STRICT_OPERANCDS = "strictOperands"

var defaultTranslator *Translator

type GetOpt struct {
	WithOperands
	WithOptions

	help HelpInterface
	settings map[string]string
	operandsCount int
	commands map[string]*Command
	command *Command
	additionalOperands []string
	additionalOptions map[string]Option
	translator Translator
}

func NewGetOpt(options string, settings map[string]string) (*GetOpt, error) {
	g := new(GetOpt)

	args := os.Args

	if len(args) >= 1 {
		SETTING_COMMAND_NAME = args[0]
	} else {
		SETTING_COMMAND_NAME = "unknown"
	}

	for setting, value := range settings {
		g.Set(setting, value)
	}

	if options != "" {
		_, err := g.AddOptionString(options)

		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (g *GetOpt) Set(setting string, value string) *GetOpt {
	switch (setting) {
	case SETTING_DEFAULT_MODE:
		defaultOptionParserMode = value
	default:
		g.settings[setting] = value
	}

	return g
}

func (g *GetOpt) Get(setting string) string {
	return g.settings[setting]
}

func (g *GetOpt) Process(args ...string) error {
	var arguments *Arguments
	if len(args) == 0 {
		args = os.Args
		arguments = NewArguments(args)
	} else {
		arguments = NewArguments(args)
	}

	setOption := func(g *GetOpt, name string, getValue getValueFunc) error {
		option := g.GetOption(name)

		if option == nil {
			if len(g.Get(SETTING_STRICT_OPTIONS)) > 0 {
				value := getValue()
				if len(value) == 0 {
					value = "1"
				}
				opt, isset := g.additionalOptions[name]
				v := opt.GetValue()
				if isset && isInt(value) && len(v) > 0 && isInt(v[0]) {
					x, _ := strconv.Atoi(value)
					y, _ := strconv.Atoi(v[0])
					value = strconv.Itoa(x + y)
				}

				newOption := Option{value: value}
				g.additionalOptions[name] = newOption
				return nil
			}

			return errors.New(fmt.Sprintf(Translate("option-unknown"), name))
		}

		return nil
	}

	setCommand := func(g *GetOpt, command *Command) error {
		g.AddOptions(command.GetOptions())
		_, err := g.AddOperands(command.GetOperands())

		if err != nil {
			return err
		}

		g.command = command

		return nil
	}

	addOperand := func(g *GetOpt, value string) error {
		operand := g.NextOperand()
		if operand != nil {
			operand.SetValue(value)
		} else if (g.Get(SETTING_STRICT_OPTIONS) != "") {
			return errors.New(fmt.Sprintf(Translate("no-more-operands"), value))
		} else {
			g.additionalOperands = append(g.additionalOperands, value)
		}

		return nil
	}

	g.additionalOptions = make(map[string]Option)
	g.additionalOperands = make([]string,0)
	g.operandsCount = 0

	arguments.Process(g, setOption, setCommand, addOperand)

	operand := g.NextOperand()

	if operand != nil && operand.IsRequired() && (!operand.IsMultiple() || len(operand.GetValue()) == 0) {
		return errors.New(fmt.Sprintf(Translate("operand-missing"), operand.GetName()))
	}

	return nil
}

func (g *GetOpt) GetOptionValue(name string) []string {
	option := g.GetOption(name)
	if option != nil {
		return option.GetValue()
	}

	return nil
}

func (g *GetOpt) GetOptions() map[string][]string {
	result := map[string][]string{}

	for _, option := range g.options {
		value := option.GetValue()
		if value != nil {
			key := string(option.GetShort())
			if key == "" {
				key = option.GetLong()
			}
			result[key] = value
			if option.GetShort() != '\x00' {
				result[string(option.GetShort())] = value
			}
			if option.GetLong() != "" {
				result[option.GetLong()] = value
			}
		}
	}

	for key, option := range g.additionalOptions {
		value := option.GetValue()
		if value != nil {
			result[key] = value
		}
	}

	return result
}

func (g *GetOpt) AddCommands(commands []*Command) *GetOpt {
	for _, command := range commands {
		g.AddCommand(command)
	}

	return g
}

func (g *GetOpt) AddCommand(command *Command) (*GetOpt, error) {
	for _, option := range command.GetOptions() {
		if g.Conflicts(option) {
			return nil, errors.New(fmt.Sprintf("%s has conflicting options", command))
		}
	}
	g.commands[command.GetName()] = command

	return g, nil
}

func (g *GetOpt) GetCommand(name string) *Command {
	if len(name) > 0 {
		v, isset := g.commands[name]
		if isset {
			return v
		}
		return nil
	}

	return g.command;
}

func (g *GetOpt) GetCommands() map[string]*Command {
	return g.commands;
}

func (g *GetOpt) HasCommands() bool {
	return len(g.commands) > 0
}

func (g *GetOpt) NextOperand() *Operand {

	if g.operandsCount <= len(g.operands) {
		operand := g.operands[g.operandsCount]
		if !operand.IsMultiple() {
			g.operandsCount++
		}
		return operand
	}

	return nil
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

func (g *GetOpt) GetOperand(index string) []string {
	operand := g.WithOperands.GetOperand(index)
	if operand != nil {
		return operand.GetValue()
	} else if (isInt(index)) {
		v, _ := strconv.Atoi(index)
		i := v - len(g.operands)
		if i >= 0 && i < len(g.additionalOperands)  {
			return []string{ g.additionalOperands[i] }
		}
	}

	return nil
}

func (g *GetOpt) SetHelp(help HelpInterface) *GetOpt {
	g.help = help
	return g
}

func (g *GetOpt) SetHelpLang(language string) bool {
	return g.SetLang(language)
}

func (g *GetOpt) SetLang(language string) bool {
	return GetTranslator().SetLanguage(language)
}

func (g *GetOpt) GetHelp() HelpInterface {
	if g.help == nil {
		g.help = NewHelp(nil)
	}

	return g.help
}

func (g *GetOpt) GetHelpText(data map[string]string) string {
	return g.GetHelp().Render(g, data)
}

func (g *GetOpt) SetScriptName(scriptName string) *GetOpt {
	return g.Set(SETTING_COMMAND_NAME, scriptName)
}

func (g *GetOpt) Parse(arguments string) {
	_ = g.Process(arguments)
}

func (g *GetOpt) Iter() map[string][]string {
	result := make(map[string][]string)

	for _, option := range g.options {
		value := option.GetValue()
		if value != nil {
			name := string(option.GetShort())
			if option.GetLong() != "" {
				name = option.GetLong()
			}
			result[name] = value
		}
	}

	return result
}

func (g *GetOpt) Count() int {
	return len(g.Iter())
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