package commands

import (
	"flag"
)

type Command interface {
	Flagset(flagset *flag.FlagSet) *flag.FlagSet
	Execute(flagset *flag.FlagSet) error
	Help()
}
