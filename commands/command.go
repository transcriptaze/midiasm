package commands

import (
	"flag"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type Command interface {
	Flagset() *flag.FlagSet
	Execute() error
	Help()

	MiddleC() lib.MiddleC
	Debug() bool
	Verbose() bool

	Config() string
}
