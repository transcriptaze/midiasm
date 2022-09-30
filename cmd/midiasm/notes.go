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
	conf      string
	out       string
	split     bool
	middleC   types.MiddleC
	transpose int
	json      bool
	verbose   bool
	debug     bool
}

func (n Notes) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("notes", flag.ExitOnError)

	flagset.StringVar(&n.out, "out", "", "Output file path")
	flagset.Var(&n.middleC, "middle-c", "Middle C convention (C3 or C4). Defaults to C3")
	flagset.IntVar(&n.transpose, "transpose", 0, "Transpose notes up or down")
	flagset.BoolVar(&n.json, "json", false, "Formats the output as JSON")
	flagset.BoolVar(&n.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&n.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (n Notes) config() string {
	return n.conf
}

func (n Notes) MiddleC() types.MiddleC {
	return n.middleC
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
