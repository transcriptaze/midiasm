package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/processors"
	"os"
)

type Notes struct {
}

func (n *Notes) Execute(smf *midi.SMF) {
	w := os.Stdout
	err := error(nil)

	if options.out != "" {
		w, err = os.Create(options.out)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = options.verbose
	eventlog.EventLog.Debug = options.debug

	p := processors.Notes{w}
	err = p.Execute(smf)
	if err != nil {
		fmt.Printf("Error %v extracting notes\n", err)
	}

}

func (n *Notes) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("notes", flag.ExitOnError)

	flagset.StringVar(&options.out, "out", "", "Output file path")
	flagset.BoolVar(&options.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&options.debug, "debug", false, "Enable debugging information")

	return flagset
}
