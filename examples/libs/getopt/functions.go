package getopt

type addOperandFunc func(*GetOpt, string) error
type getValueFunc func() string
type handlerFunc func()
type setOptionFunc func(*GetOpt, string, getValueFunc) error
type setCommandFunc func(*GetOpt, *Command) error
type validationFunc func(args ...string) bool
