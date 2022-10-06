package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/transcriptaze/midiasm/ops/assemble"
)

type Assemble struct {
	command
	out string
}

var ASSEMBLE = Assemble{
	command: command{
		c4:      false,
		verbose: false,
		debug:   false,
	},
}

func (a *Assemble) Flagset() *flag.FlagSet {
	flagset := a.command.flagset("assemble")

	flagset.StringVar(&a.out, "out", "", "Output file path")

	a.flags = flagset

	return flagset
}

func (a Assemble) Help() {
	fmt.Println()
	fmt.Println("  Assembles a MIDI file from a text or JSON source.")
	fmt.Println()
	fmt.Println("    midiasm assemble [--debug] [--verbose] [--C4] [--out <MIDI file>] <file>")
	fmt.Println()
	fmt.Println("      --out <file>  Output MIDI file. Default is to use the input file name with a .midi extension.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug    Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose  Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm assemble --debug --verbose --out one-time.midi one-time.txt")
	fmt.Println()
}

func (a Assemble) Execute() error {
	filename := a.flags.Arg(0)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var assembler assemble.Assembler

	switch filepath.Ext(filename) {
	case ".json":
		assembler = assemble.NewJSONAssembler()

	default:
		assembler = assemble.NewTextAssembler()
	}

	if midi, err := assembler.Assemble(bytes); err != nil {
		return err
	} else {
		fmt.Printf(">>>>>> MIDI: %v bytes", len(midi))
	}

	return nil
}
