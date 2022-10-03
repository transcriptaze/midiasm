package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/ops/transpose"
)

type Transpose struct {
	command
	out       string
	semitones int
	flags     *flag.FlagSet
}

var TRANSPOSE = Transpose{}

func (t *Transpose) Flagset() *flag.FlagSet {
	flagset := t.command.flagset("transpose")

	flagset.StringVar(&t.out, "out", "", "Output file path")
	flagset.IntVar(&t.semitones, "semitones", 0, "Number of semitones to transpose notes (+ve is up, -ve is down")

	t.flags = flagset

	return flagset
}

func (t Transpose) Help() {
	fmt.Println()
	fmt.Println("  Transposes the key of the notes (and key signature) and writes it back as MIDI file.")
	fmt.Println()
	fmt.Println("    midiasm transpose [--debug] [--verbose] [--C4] --semitones <steps> --out <file> <MIDI file>")
	fmt.Println()
	fmt.Println("      --semitones <N>  Number of semitones to transpose up or down. Defaults to 0.")
	fmt.Println("      --out <file>     (required) Destination file for the transposed MIDI. ")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug    Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose  Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm transpose --debug --verbose --semitones +5 --out one-time+5.mid one-time.mid")
	fmt.Println()
}

func (t Transpose) Execute() error {
	filename := t.flags.Arg(0)

	smf, err := t.decode(filename)
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

	return t.execute(smf)
}

func (t Transpose) execute(smf *midi.SMF) error {
	op := transpose.Transpose{}

	transposed, err := op.Execute(smf, t.semitones)
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
