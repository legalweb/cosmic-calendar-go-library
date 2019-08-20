package getopt

type addOperandFunc func(GetOpt, string)
type getValueFunc func() string
type setOptionFunc func(GetOpt, string, getValueFunc)
type setCommandFunc func(GetOpt)
type validationFunc func(args ...string) bool
