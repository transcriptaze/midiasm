package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/midi/types"
	"github.com/transcriptaze/midiasm/ops/notes"
)

type Notes struct {
	command
	out       string
	transpose int
	json      bool
}

var NOTES = Notes{
	command: command{
		middleC: types.C3,
	},
}

func (n *Notes) flagset() *flag.FlagSet {
	flagset := n.command.flagset("notes")

	flagset.StringVar(&n.out, "out", "", "Output file path")
	flagset.IntVar(&n.transpose, "transpose", 0, "Transpose notes up or down")
	flagset.BoolVar(&n.json, "json", false, "Formats the output as JSON")

	return flagset
}

func (n Notes) Execute(smf *midi.SMF) {
	w := os.Stdout
	err := error(nil)

	if n.out != "" {
		w, err = os.Create(n.out)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = n.verbose
	eventlog.EventLog.Debug = n.debug

	p := notes.Notes{
		JSON:      n.json,
		Transpose: n.transpose,
		Writer:    w,
	}

	err = p.Execute(smf)
	if err != nil {
		fmt.Printf("Error %v extracting notes\n", err)
	}
}
