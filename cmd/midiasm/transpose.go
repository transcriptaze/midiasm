package main

import (
	"flag"
	"fmt"
	"os"

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

	if bytes, err := t.read(filename); err != nil {
		return err
	} else if t.out == "" {
		return fmt.Errorf("missing 'out' file")
	} else {
		return t.execute(bytes)
	}
}

func (t Transpose) execute(bytes []byte) error {
	op := transpose.Transpose{}

	transposed, err := op.Execute(bytes)
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
