package command

import "github.com/spf13/pflag"

type FlagSetter interface {
	MarkFlagRequired(name string) error
	Flags() *pflag.FlagSet
}
