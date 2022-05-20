package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/ops/notes"
	"os"
)

type Notes struct {
	conf    string
	out     string
	split   bool
	verbose bool
	debug   bool
}

func (n *Notes) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("notes", flag.ExitOnError)

	flagset.StringVar(&n.out, "out", "", "Output file path")
	flagset.BoolVar(&n.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&n.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (n *Notes) config() string {
	return n.conf
}

func (n *Notes) Execute(smf *midi.SMF) {
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

	p := notes.Notes{w}
	err = p.Execute(smf)
	if err != nil {
		fmt.Printf("Error %v extracting notes\n", err)
	}
}
