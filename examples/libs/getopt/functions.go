package getopt

type addOperandFunc func(string) error
type getValueFunc func() string
type handlerFunc func(opt *GetOpt)
type setOptionFunc func(string, getValueFunc) error
type setCommandFunc func(*Command) error
type validationFunc func(...string) bool
