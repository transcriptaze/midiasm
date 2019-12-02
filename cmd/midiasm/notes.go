package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"os"
)

func notes(smf *midi.SMF) {
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

	smf.Notes(w)
}

func notesFlagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("notes", flag.ExitOnError)

	flagset.StringVar(&options.out, "out", "", "Output file path")
	flagset.BoolVar(&options.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&options.debug, "debug", false, "Enable debugging information")

	return flagset
}
