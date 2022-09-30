package main

import (
	"flag"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/types"
)

type Command interface {
	flagset() *flag.FlagSet
	config() string
	MiddleC() types.MiddleC
	Execute(*midi.SMF)
}

type command struct {
	conf    string
	middleC types.MiddleC
	verbose bool
	debug   bool
}

func (c command) config() string {
	return c.conf
}

func (c command) MiddleC() types.MiddleC {
	return c.middleC
}

func (c *command) flagset(name string) *flag.FlagSet {
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	flagset.Var(&c.middleC, "middle-c", "Middle C convention (C3 or C4). Defaults to C3")
	flagset.BoolVar(&c.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&c.debug, "debug", false, "Enable debugging information")

	return flagset
}
