package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	impl "github.com/transcriptaze/midiasm/ops/export"
)

type export struct {
	out string
}

var Export = export{}

func (x *export) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.StringVar(&x.out, "out", "", "Output file path (or directory for split files)")

	return flagset
}

func (x export) Help() {
	fmt.Println()
	fmt.Println("  Extracts the MIDI information as JSON for use with other tools (e.g. jq).")
	fmt.Println()
	fmt.Println("    midiasm export [--debug] [--verbose] [--C4] [--out <file>] <MIDI file>")
	fmt.Println()
	fmt.Println("      <MIDI file>  MIDI file to export as JSON.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --out <file>  Writes the JSON to a file. Default is to write to stdout.")
	fmt.Println("      --C4          Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug       Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose     Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm export --debug --verbose --out one-time.json one-time.mid")
	fmt.Println()
}

func (x export) Execute(flagset *flag.FlagSet) error {
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

	return x.execute(smf)
}

func (x export) execute(smf *midi.SMF) error {
	op, err := impl.NewExport()
	if err != nil {
		return err
	}

	return x.write(op, smf)
}

func (x export) write(op *impl.Export, smf *midi.SMF) error {
	out := os.Stdout

	if x.out != "" {
		w, err := os.Create(x.out)
		if err != nil {
			return err
		}

		defer w.Close()

		out = w
	}

	return op.Export(smf, out)
}
