package main

import (
	"flag"
	"fmt"

	"github.com/transcriptaze/midiasm/commands"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type tsv struct {
	out     string
	conf    string
	c4      bool
	verbose bool
	debug   bool
}

var TSV tsv

func (t tsv) GetCommand() (string, commands.Command) {
	return "tsv", TSV
}

func (t tsv) Flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("tsv", flag.ExitOnError)

	flagset.StringVar(&t.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&t.c4, "C4", t.c4, "Sets middle C to C4 (Yamaho convention). Defaults to C3")
	flagset.BoolVar(&t.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&t.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (t tsv) Execute() error {
	return fmt.Errorf("NOT IMPLEMENTED")
}

func (t tsv) Help() {
	fmt.Println()
	fmt.Println("  Extracts the MIDI information as TSV for use with e.g. a spreadsheet.")
	fmt.Println()
	fmt.Println("    midiasm tsv [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>")
	fmt.Println()
	fmt.Println("      <MIDI file>  MIDI file to export as JSON.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --out <file>  Writes the TSV to a file. Default is to write to stdout.")
	fmt.Println("      --C4          Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug       Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose     Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm tsv --debug --verbose --out one-time.tsv one-time.mid")
	fmt.Println()
}

func (t tsv) MiddleC() lib.MiddleC {
	if t.c4 {
		return lib.C4
	}

	return lib.C3
}

func (t tsv) Config() string {
	return t.conf
}

func (t tsv) Debug() bool {
	return t.debug
}

func (t tsv) Verbose() bool {
	return t.verbose
}
