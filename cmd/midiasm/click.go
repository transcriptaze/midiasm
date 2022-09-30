package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/midi/types"
	"github.com/transcriptaze/midiasm/ops/click"
)

type Click struct {
	conf    string
	out     string
	middleC types.MiddleC
	split   bool
	verbose bool
	debug   bool
}

func (c *Click) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("notes", flag.ExitOnError)

	flagset.StringVar(&c.out, "out", "", "Output file path")
	flagset.Var(&c.middleC, "middle-c", "Middle C convention (C3 or C4). Defaults to C3")
	flagset.BoolVar(&c.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&c.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (c *Click) config() string {
	return c.conf
}

func (c *Click) MiddleC() types.MiddleC {
	return c.middleC
}

func (c *Click) Execute(smf *midi.SMF) {
	var w = os.Stdout
	var err error

	if c.out != "" {
		if w, err = os.Create(c.out); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = c.verbose
	eventlog.EventLog.Debug = c.debug

	p := click.ClickTrack{w}
	if err = p.Execute(smf); err != nil {
		fmt.Printf("Error %v extracting click track\n", err)
	}
}
