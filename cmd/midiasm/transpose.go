package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/midi/types"
	"github.com/transcriptaze/midiasm/ops/transpose"
)

type Transpose struct {
	command
	out       string
	transpose int
}

var TRANSPOSE = Transpose{
	command: command{
		middleC: types.C3,
	},
}

func (t *Transpose) flagset() *flag.FlagSet {
	flagset := t.command.flagset("transpose")

	flagset.StringVar(&t.out, "out", "", "Output file path")
	flagset.IntVar(&t.transpose, "transpose", 0, "Transpose notes up or down")

	return flagset
}

func (t Transpose) Execute(filename string) error {
	eventlog.EventLog.Verbose = t.verbose
	eventlog.EventLog.Debug = t.debug

	smf, err := t.decode(filename)
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

	return t.execute(smf)
}

func (t Transpose) execute(smf *midi.SMF) error {
	op := transpose.Transpose{}

	transposed, err := op.Execute(smf, t.transpose)
	if err != nil {
		return err
	}

	if t.out != "" {
		if w, err := os.Create(t.out); err != nil {
			return err
		} else {
			defer w.Close()

			if _, err := w.Write(transposed); err != nil {
				return err
			}
		}
	}

	return nil
}
