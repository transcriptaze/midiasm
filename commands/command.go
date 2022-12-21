package commands

import (
	"flag"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type Command interface {
	Flagset() *flag.FlagSet
	Execute(flagset *flag.FlagSet) error
	Help()

	MiddleC() lib.MiddleC
	Debug() bool
	Verbose() bool

	Config() string
}
