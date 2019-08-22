package getopt

import "strings"

const TEMPLATE_USAGE = "usageTemplate"
const TEMPLATE_OPTIONS = "optionsTemplate"
const TEMPLATE_COMMANDS = "commandsTemplate"
const DESCRIPTION = "description"
const MAX_WIDTH = "maxWidth"
const HIDE_OPERANDS = "hideOperands"

type HelpInterface interface {
	Render(*GetOpt, map[string]string) string
}

type Help struct {
	usageTemplate string
	optionsTemplate string
	commandsTemplate string
	settings map[string]string
	texts map[string]string
	getOpt *GetOpt
	screenWidth int
}

func NewHelp(settings map[string]string) *Help {
	h := new(Help)

	h.settings[MAX_WIDTH] = "120"
	h.texts["placeholder"] = "<>"
	h.texts["optional"] = "[]"
	h.texts["multiple"] = "..."
	h.texts["options-listing"] = ", "

	for setting, value := range settings {
		h.Set(setting, value)
	}

	return h
}

func (h *Help) Set(setting string, value string) *Help {
	switch(setting) {
	case "optionsTemplate":
		h.SetOptionsTemplate(value)
	case "commandsTemplate":
		h.SetCommandsTemplate(value)
	case "usageTemplate":
		h.SetUsageTemplate(value)
	default:
		h.settings[setting] = value
	}

	return h
}

func (h *Help) SetTexts(texts map[string]string) *Help {
	for k, v := range texts {
		h.texts[k] = v
	}

	for k, v := range h.texts {
		h.texts[k] = strings.ReplaceAll(v, "\r\n", "\n")
	}

	return h
}

func (h *Help) GetText(key string) string {
	v, isset := h.texts[key]
	if !isset {
		return Translate(key)
	}

	return v
}

func (h *Help) Render(getopt *GetOpt, data map[string]string) string {
	helpText := ""

	h.getOpt = getopt

	for setting, value := range data {
		h.Set(setting, value)
	}

	if h.usageTemplate != "" {
		helpText = h.renderTemplate(h.usageTemplate, getopt, getopt.getCommand())
	} else {
		helpText = h.renderUsage()
	}

	_, isset := h.settings[HIDE_OPERANDS]

	if getopt.HasOperands() && isset {
		helpText += h.renderOperands()
	}

	if getopt.HasOptions() {
		if h.optionsTemplate != "" {
			helpText += h.renderTemplate(h.optionsTemplate, getopt.WithOptions.GetOptions())
		} else {
			helpText += h.renderOptions()
		}
	}

	if getopt.GetCommand() && getopt.HasCommands() {
		if h.commandsTemplate != "" {
			helpText += h.renderTemplate(h.commandsTemplate, getopt.GetCommands())
		} else {
			helpText += h.RenderCommands()
		}
	}

	return helpText
}

func (h *Help) surround(text string, with string) string {
	return string(with[0]) + text + with[:-1]
}

func (h *Help) renderUsage() string {
	return h.GetText("usage-title") +
		h.getOpt.get(SETTING_SCRIPT_NAME) +
		" " +
		h.renderUsageCommand() +
		h.renderUsageOptions() +
		h.renderUsageOperands() +
		"\n\n" +
		h.renderDescription()
}

func (h *Help) renderOperands() string {
	data := make([]string, 0)
	definitionWidth := 0
	hasDescriptions := false

	for _, operand := range h.getOpt.WithOperands.GetOperands() {
		definition := h.surround(operand.GetName(), h.texts["placeholder"])
		if !operand.IsRequired() {
			definition = h.surround(definition, h.texts["optional"])
		}

		if len(definition) > definitionWidth {
			definitionWidth = len(definition)
		}

		if operand.GetDescription() != "" {
			hasDescriptions = true
		}

		data = append(data, operand.GetDescription())
	}

	if hasDescriptions {
		return ""
	}

	return h.GetText("operands-title") + h.renderColumns(definitionWidth, data) + "\n"
}

func (h *Help) renderOptions() string {
	data := make([]string, 0)
	definitionWidth := 0

	for _, option := range h.getOpt.WithOptions.GetOptions() {
		optionStrings := make([]string,0)
		if option.GetShort() != '\x00' {
			optionStrings = append(optionStrings, "-" + string(option.GetShort()))
		}
		if option.GetLong() != "" {
			optionStrings = append(optionStrings, "--" + string(option.GetLong()))
		}

		definition := strings.Join(optionStrings,h.texts["options-listing"])

		if option.GetMode() != NO_ARGUMENT {
			argument := h.surround(option.GetArgument().GetName(), h.texts["placeholder"])
			if option.GetMode() == OPTIONAL_ARGUMENT {
				argument = h.surround(argument, h.texts["optional"])
			}

			definition += " " + argument
		}

		if len(definition) > definitionWidth {
			definitionWidth = len(definition)
		}

		data = append(data, definition, option.GetDescription())
	}

	return h.GetText("options-title") + h.renderColumns(definitionWidth, data) + "\n"
}

func (h *Help) renderCommands() {
	data := make([]string, 0)
	nameWidth := 0

	for _, command := range h.getOpt.GetCommands() {
		if len(command.GetName()) > nameWidth {
			nameWidth = len(command.GetName())
		}

		data = append(data, command.GetName(), command.GetShortDescription())
	}

	return h.GetText("commands-title") + h.renderColumns(nameWidth, data) + "\n"
}

func (h *Help) renderUsageCommand() string {
	command := h.getOpt.GetCommand("")

	if command != nil {
		return command.GetName() + " "
	} else if h.getOpt.HasCommands() {
		return h.surround(h.GetText("usage-command"), h.texts["placeholder"]) + " "
	}

	return ""
}

func (h *Help) renderUsageOptions() string {
	if h.getOpt.HasOptions() || len(h.getOpt.Get(SETTING_STRICT_OPTIONS)) > 0 {
		return h.surround(h.GetText("usage-options"), h.texts["optional"]) + " "
	}

	return ""
}

