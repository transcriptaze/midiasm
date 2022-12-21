package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
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
		c4:      false,
		verbose: false,
		debug:   false,
	},
}

func (n *Notes) Flagset() *flag.FlagSet {
	flagset := n.command.flagset("notes")

	flagset.StringVar(&n.out, "out", "", "Output file path")
	flagset.IntVar(&n.transpose, "transpose", 0, "Transpose notes up or down")
	flagset.BoolVar(&n.json, "json", false, "Formats the output as JSON")

	return flagset
}

func (n Notes) Help() {
	fmt.Println()
	fmt.Println("  Extracts the NoteOn and NoteOff events to generate a list of notes with start times and durations.")
	fmt.Println()
	fmt.Println("    midiasm notes [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>")
	fmt.Println()
	fmt.Println("      --out <file>  Writes the notes to a file. Default is to write to stdout.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug    Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose  Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm notes --debug --verbose --out one-time.notes one-time.mid")
	fmt.Println()
}

func (n Notes) Execute(flagset *flag.FlagSet) error {
	filename := flagset.Arg(0)

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

	op := notes.Notes{
		JSON:      n.json,
		Transpose: n.transpose,
		Writer:    w,
	}

	return op.Execute(smf)
}
