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

func (c Click) Execute(filename string) error {
	smf, err := c.decode(filename)
	if err != nil {
		return err
	}

	if errors := smf.Validate(); len(errors) > 0 {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "WARNING: there are validation errors:\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "         ** %v\n", e)
		}
		fmt.Fprintln(os.Stderr)
	}

	return c.execute(smf)
}

func (c Click) execute(smf *midi.SMF) error {
	var w = os.Stdout
	var err error

	if c.out != "" {
		if w, err = os.Create(c.out); err != nil {
			return err
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = c.verbose
	eventlog.EventLog.Debug = c.debug

	track := click.ClickTrack{w}

	return track.Execute(smf)
}
