package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/ops/click"
	"os"
)

type Click struct {
	conf    string
	out     string
	split   bool
	verbose bool
	debug   bool
}

func (c *Click) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("notes", flag.ExitOnError)

	flagset.StringVar(&c.out, "out", "", "Output file path")
	flagset.BoolVar(&c.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&c.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (c *Click) config() string {
	return c.conf
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
