package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
)

type humanise struct {
	out string
}

var Humanise = humanise{}

func (h *humanise) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.StringVar(&TSV.out, "out", "", "Output file path")

	return flagset
}

func (h humanise) Help() {
	fmt.Println()
	fmt.Println("  Randomises the MIDI event timing")
	fmt.Println()
	fmt.Println("    midiasm humanise [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>")
	fmt.Println()
	fmt.Println("      <MIDI file>  MIDI file to 'humanise'.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --out <file>          Writes the TSV to a file. Default is to write to stdout.")
	fmt.Println("      --C4                  Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug               Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose             Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm humanise --debug --verbose --out gnossienes.mid")
	fmt.Println()
}

func (h humanise) Execute(flagset *flag.FlagSet) error {
	filename := flagset.Arg(0)

	if b, err := os.ReadFile(filename); err != nil {
		return err
	} else if smf, err := h.decode(bytes.NewBuffer(b)); err != nil {
		return err
	} else if err := h.validate(smf); err != nil {
		return err
	} else if v, err := h.humanise(smf); err != nil {
		return err
	} else if err := h.write(v); err != nil {
		return err
	}

	return nil
}

func (h humanise) humanise(smf *midi.SMF) (any, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED")
}

func (h humanise) write(v any) error {
	return fmt.Errorf("** NOT IMPLEMENTED **")
}

func (h humanise) decode(r io.Reader) (*midi.SMF, error) {
	decoder := midifile.NewDecoder()

	if smf, err := decoder.Decode(r); err != nil {
		return nil, err
	} else if smf == nil {
		return nil, fmt.Errorf("failed to decode MIDI file")
	} else {
		return smf, nil
	}
}

func (h humanise) validate(smf *midi.SMF) error {
	errors := smf.Validate()

	if len(errors) > 0 {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "WARNING: there are validation errors:\n")
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "         ** %v\n", e)
		}
		fmt.Fprintln(os.Stderr)
	}

	return nil
}
