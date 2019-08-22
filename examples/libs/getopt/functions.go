package getopt

type addOperandFunc func(string) error
type getValueFunc func() string
type handlerFunc func()
type setOptionFunc func(string, getValueFunc) error
type setCommandFunc func(*Command) error
type validationFunc func(args ...string) bool
