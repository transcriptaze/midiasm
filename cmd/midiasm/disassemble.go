package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/ops/disassemble"
)

type Disassemble struct {
	out       string
	split     bool
	templates string
}

var DISASSEMBLE = Disassemble{}

func (d *Disassemble) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	flagset.StringVar(&d.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&d.split, "split", false, "Create separate file for each track. Defaults to the same directory as the MIDI file.")
	flagset.StringVar(&d.templates, "templates", "", "Loads the formatting templates from a file")

	return flagset
}

func (d Disassemble) Help() {
	fmt.Println()
	fmt.Println("  Disassembles a MIDI file and displays the tracks in a human readable format.")
	fmt.Println()
	fmt.Println("    midiasm [--debug] [--verbose] [--C4] [--split] [--out <file>] <MIDI file>")
	fmt.Println()
	fmt.Println("      --out <file>  Writes the disassembly to a file. Default is to write to stdout.")
	fmt.Println("      --split       Writes each track to a separate file. Default is `false`.")
	fmt.Println()
	fmt.Println("    Options:")
	fmt.Println()
	fmt.Println("      --C4       Uses C4 as middle C (Yamaha convention). Defaults to C3.")
	fmt.Println("      --debug    Displays internal information while processing a MIDI file. Defaults to false")
	fmt.Println("      --verbose  Enables 'verbose' logging. Defaults to false")
	fmt.Println()
	fmt.Println("    Example:")
	fmt.Println()
	fmt.Println("      midiasm --debug --verbose --out one-time.txt one-time.mid")
	fmt.Println()
}

func (p Disassemble) Execute(flagset *flag.FlagSet) error {
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

	op, err := disassemble.NewDisassemble()
	if err != nil {
		return err
	}

	if p.templates != "" {
		f, err := os.Open(p.templates)
		if err == nil {
			err = op.LoadTemplates(f)
			f.Close()
		}

		if err != nil {
			return err
		}
	}

	if p.split && p.out == "" {
		return p.separate(op, smf, filename)
	} else if p.split && p.out != "" {
		return p.separate(op, smf, p.out)
	} else {
		return p.write(op, smf)
	}
}

func (p Disassemble) write(op *disassemble.Disassemble, smf *midi.SMF) error {
	out := os.Stdout

	if p.out != "" {
		w, err := os.Create(p.out)
		if err != nil {
			return err
		}

		defer w.Close()

		out = w
	}

	return op.Print(smf, "document", out)
}

func (p Disassemble) separate(op *disassemble.Disassemble, smf *midi.SMF, file string) error {
	// Get base filename and Create out directory
	base := strings.TrimSuffix(path.Base(file), path.Ext(file))
	dir := path.Dir(file)

	if p.out != "" {
		dir = p.out
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return err
		}
	}

	// Print MThd
	filename := fmt.Sprintf("%s.MThd", base)
	w, err := os.Create(path.Join(dir, filename))
	if err != nil {
		return err
	}

	if err = op.Print(smf.MThd, "MThd", w); err != nil {
		fmt.Printf("Error %v extracting MThd information\n", err)
	}

	w.Close()

	// Print tracks
	for _, mtrk := range smf.Tracks {
		filename := fmt.Sprintf("%s-%d.MTrk", base, mtrk.TrackNumber)
		w, err := os.Create(path.Join(dir, filename))
		if err != nil {
			return err
		}

		if err = op.Print(mtrk, "MTrk", w); err != nil {
			fmt.Printf("Error %v extracting MTrk%d information\n", mtrk.TrackNumber, err)
		}

		w.Close()
	}

	return nil
}
