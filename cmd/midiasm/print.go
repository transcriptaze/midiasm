package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/midi/types"
	"github.com/transcriptaze/midiasm/ops/print"
)

type Print struct {
	command
	out       string
	split     bool
	templates string
}

var PRINT = Print{
	command: command{
		middleC: types.C3,
	},
}

func (p *Print) flagset() *flag.FlagSet {
	flagset := p.command.flagset("print")

	flagset.StringVar(&p.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&p.split, "split", false, "Create separate file for each track. Defaults to the same directory as the MIDI file.")
	flagset.StringVar(&p.templates, "templates", "", "Loads the formatting templates from a file")

	return flagset
}

func (p Print) Execute(filename string) error {
	smf, err := p.decode(filename)
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

	return p.execute(smf)
}

func (p Print) execute(smf *midi.SMF) error {
	eventlog.EventLog.Verbose = p.verbose
	eventlog.EventLog.Debug = p.debug

	op, err := print.NewPrint()
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

	if p.split {
		return p.separate(op, smf)
	} else {
		return p.write(op, smf)
	}
}

func (p Print) write(op *print.Print, smf *midi.SMF) error {
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

func (p Print) separate(op *print.Print, smf *midi.SMF) error {
	// Get base filename and Create out directory
	base := strings.TrimSuffix(path.Base(smf.File), path.Ext(smf.File))
	dir := path.Dir(smf.File)

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
