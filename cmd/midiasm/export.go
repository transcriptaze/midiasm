package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/midi/types"
	"github.com/transcriptaze/midiasm/ops/export"
)

type Export struct {
	command
	out string
}

var EXPORT = Export{
	command: command{
		middleC: types.C3,
	},
}

func (x *Export) flagset() *flag.FlagSet {
	flagset := x.command.flagset("export")

	flagset.StringVar(&x.out, "out", "", "Output file path (or directory for split files)")

	return flagset
}

func (x Export) Execute(smf *midi.SMF) {
	eventlog.EventLog.Verbose = x.verbose
	eventlog.EventLog.Debug = x.debug

	op, err := export.NewExport()
	if err != nil {
		fmt.Printf("ERROR  %v\n", err)
		return
	}

	x.write(op, smf)
}

func (x Export) write(op *export.Export, smf *midi.SMF) {
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
