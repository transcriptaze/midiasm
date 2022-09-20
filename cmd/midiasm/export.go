package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/ops/export"
)

type Export struct {
	conf    string
	out     string
	verbose bool
	debug   bool
}

func (x *Export) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("print", flag.ExitOnError)

	flagset.StringVar(&x.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&x.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&x.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (x *Export) config() string {
	return x.conf
}

func (x *Export) Execute(smf *midi.SMF) {
	eventlog.EventLog.Verbose = x.verbose
	eventlog.EventLog.Debug = x.debug

	op, err := export.NewExport()
	if err != nil {
		fmt.Printf("ERROR  %v\n", err)
		return
	}

	x.write(op, smf)
}

func (x *Export) write(op *export.Export, smf *midi.SMF) {
	out := os.Stdout

	if x.out != "" {
		w, err := os.Create(x.out)
		if err != nil {
			fmt.Printf("ERROR  %v\n", err)
			return
		}

		defer w.Close()

		out = w
	}

	err := op.Export(smf, out)
	if err != nil {
		fmt.Printf("ERROR  %v\n", err)
	}
}
