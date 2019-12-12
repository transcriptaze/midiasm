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
}

func (p *Print) Execute(smf *midi.SMF) {
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

	x := processors.Print{w}
	err = x.Execute(smf)
	if err != nil {
		fmt.Printf("Error %v extracting notes\n", err)
	}
}

func (p *Print) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("print", flag.ExitOnError)

	flagset.StringVar(&options.out, "out", "", "Output file path")
	flagset.BoolVar(&options.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&options.debug, "debug", false, "Enable debugging information")

	return flagset
}
