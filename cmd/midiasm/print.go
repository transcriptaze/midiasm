package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/processors"
	"os"
)

type Print struct {
	out     string
	split   bool
	verbose bool
	debug   bool
}

func (p *Print) Execute(smf *midi.SMF) {
	w := os.Stdout
	err := error(nil)

	if p.out != "" {
		w, err = os.Create(p.out)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = p.verbose
	eventlog.EventLog.Debug = p.debug

	x := processors.Print{w}
	err = x.Execute(smf)
	if err != nil {
		fmt.Printf("Error %v extracting notes\n", err)
	}
}

func (p *Print) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("print", flag.ExitOnError)

	flagset.StringVar(&p.out, "out", "", "Output file path")
	flagset.BoolVar(&p.split, "split", false, "Create separate file for each track")
	flagset.BoolVar(&p.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&p.debug, "debug", false, "Enable debugging information")

	return flagset
}
