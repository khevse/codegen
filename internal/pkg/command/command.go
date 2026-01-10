package command

type Command interface {
	Name() string
	ShortName() string
	InitFlags(flagSetter FlagSetter) error
	Execute() error
}
