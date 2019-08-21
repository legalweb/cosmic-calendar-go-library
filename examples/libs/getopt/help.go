package getopt

import "strings"

const TEMPLATE_USAGE = "usageTemplate"
const TEMPLATE_OPTIONS = "optionsTemplate"
const TEMPLATE_COMMANDS = "commandsTemplate"
const DESCRIPTION = "description"
const MAX_WIDTH = "maxWidth"
const HIDE_OPERANDS = "hideOperands"

type HelpInterface interface {
	render(GetOpt, []string) string
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

func (h *Help) renderOperands() {
	data := make(map[string]string)
	definitionWidth := 0
	hasDescriptions := false

	for _, operand := range h.getOpt.GetOperandObjects() {
		definition := h.surround(operand.GetName(), h.texts["placeholder"])
		if !operand.Is
	}
}