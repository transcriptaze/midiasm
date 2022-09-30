package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/midi/types"
)

type Assemble struct {
	conf    string
	out     string
	middleC types.MiddleC
	split   bool
	verbose bool
	debug   bool
}

func (a Assemble) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("assemble", flag.ExitOnError)

	flagset.StringVar(&a.out, "out", "", "Output file path")
	flagset.Var(&a.middleC, "middle-c", "Middle C convention (C3 or C4). Defaults to C3")
	flagset.BoolVar(&a.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&a.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (a Assemble) config() string {
	return a.conf
}

func (a Assemble) MiddleC() types.MiddleC {
	return a.middleC
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
