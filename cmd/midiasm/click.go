package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/ops/click"
)

type Click struct {
	out string
}

var CLICK = Click{}

func (c *Click) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.StringVar(&c.out, "out", "", "Output file path")

	return flagset
}

func (c Click) Help() {
	fmt.Println()
	fmt.Println("  Extracts the _beats_ from the MIDI file in a format that can be used to create a click track.")
	fmt.Println()
	fmt.Println("    midiasm click [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>`")
	fmt.Println()
	fmt.Println("      --out <file>  Writes the click track to a file. Default is to write to stdout.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug    Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose  Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm click --debug --verbose --out one-time.click one-time.mid")
	fmt.Println()
}

func (c Click) Execute(flagset *flag.FlagSet) error {
	filename := flagset.Arg(0)

	smf, err := decode(filename)
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

	return c.execute(smf)
}

func (c Click) execute(smf *midi.SMF) error {
	var w = os.Stdout
	var err error

	if c.out != "" {
		if w, err = os.Create(c.out); err != nil {
			return err
		}

		defer w.Close()
	}

	track := click.ClickTrack{w}

	return track.Execute(smf)
}
