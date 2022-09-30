package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
)

type Assemble struct {
	command
	out string
}

func (a *Assemble) flagset() *flag.FlagSet {
	flagset := a.command.flagset("assemble")

	flagset.StringVar(&a.out, "out", "", "Output file path")

	return flagset
}

func (a Assemble) Execute(smf *midi.SMF) {
	var w = os.Stdout
	var err error

	if a.out != "" {
		if w, err = os.Create(a.out); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = a.verbose
	eventlog.EventLog.Debug = a.debug

	// p := click.ClickTrack{w}
	// if err = p.Execute(smf); err != nil {
	// 	fmt.Printf("Error %v extracting click track\n", err)
	// }
	fmt.Printf("NOT IMPLEMENTED\n")
}
