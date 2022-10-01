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

func (x Export) Execute(filename string) error {
	smf, err := x.decode(filename)
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

	return x.execute(smf)
}

func (x Export) execute(smf *midi.SMF) error {
	eventlog.EventLog.Verbose = x.verbose
	eventlog.EventLog.Debug = x.debug

	op, err := export.NewExport()
	if err != nil {
		return err
	}

	return x.write(op, smf)
}

func (x Export) write(op *export.Export, smf *midi.SMF) error {
	out := os.Stdout

	if x.out != "" {
		w, err := os.Create(x.out)
		if err != nil {
			return err
		}

		defer w.Close()

		out = w
	}

	return op.Export(smf, out)
}
