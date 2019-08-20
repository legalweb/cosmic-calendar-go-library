package getopt

import "os"

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
}

func (o *GetOpt) Process(args ...string) error {
	if len(args) == 0 {
		args = os.Args
	}

}

func (o *GetOpt) GetOperands() {
	operandsValues := make([]string, 0)

	for _, operand := range o.WithOperands.GetOperands() {
		value := operand.GetValue()

		if value == nil {
			continue
		}

		operandsValues = append(operandsValues, value...)
	}

	return append(operandsValues, o.additionalOperands)
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