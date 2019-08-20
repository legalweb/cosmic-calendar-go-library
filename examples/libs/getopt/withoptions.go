package getopt

type WithOptions struct {
	options []Option
	optionMapping map[string]Option
}

func (w *WithOptions) AddOptionString(option string) (*WithOptions, error) {
	op := new(OptionParser)
	options, err := op.ParseString(option)

	if err != nil {
		return nil, err
	}
}

func (w *WithOptions) AddOptions(options []Option) *WithOptions {
	if
}