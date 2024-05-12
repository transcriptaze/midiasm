package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/transcriptaze/midiasm/ops/assemble"
)

type Assemble struct {
	out string
}

var ASSEMBLE = Assemble{}

func (a *Assemble) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.StringVar(&a.out, "out", "", "Output file path")

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

func (a Assemble) Execute(flagset *flag.FlagSet) error {
	filename := flagset.Arg(0)

	var r io.Reader
	if b, err := os.ReadFile(filename); err != nil {
		return err
	} else {
		r = bytes.NewBuffer(b)
	}

	var assembler assemble.Assembler

	switch filepath.Ext(filename) {
	case ".json":
		assembler = assemble.NewJSONAssembler()

	default:
		assembler = assemble.NewTextAssembler()
	}

	midi, err := assembler.Assemble(r)
	if err != nil {
		return err
	}

	if a.out != "" {
		if err := os.WriteFile(a.out, midi, 0660); err != nil {
			return err
		}
	}

	return nil
}
