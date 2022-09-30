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
	command
	out string
}

var CLICK = Click{
	command: command{
		middleC: types.C3,
	},
}

func (c *Click) flagset() *flag.FlagSet {
	flagset := c.command.flagset("click")

	flagset.StringVar(&c.out, "out", "", "Output file path")

	return flagset
}

func (c Click) Execute(smf *midi.SMF) {
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
