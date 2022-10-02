package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/ops/notes"
)

type Notes struct {
	command
	out       string
	transpose int
	json      bool
}

var NOTES = Notes{}

func (n *Notes) flagset() *flag.FlagSet {
	flagset := n.command.flagset("notes")

	flagset.StringVar(&n.out, "out", "", "Output file path")
	flagset.IntVar(&n.transpose, "transpose", 0, "Transpose notes up or down")
	flagset.BoolVar(&n.json, "json", false, "Formats the output as JSON")

	return flagset
}

func (n Notes) Execute(filename string) error {
	smf, err := n.decode(filename)
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

	return n.execute(smf)
}

func (n Notes) execute(smf *midi.SMF) error {
	w := os.Stdout
	err := error(nil)

	if n.out != "" {
		w, err = os.Create(n.out)
		if err != nil {
			return err
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = n.verbose
	eventlog.EventLog.Debug = n.debug

	op := notes.Notes{
		JSON:      n.json,
		Transpose: n.transpose,
		Writer:    w,
	}

	return op.Execute(smf)
}
