package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/eventlog"
	"github.com/transcriptaze/midiasm/ops/print"
)

type Print struct {
	conf      string
	out       string
	split     bool
	templates string
	verbose   bool
	debug     bool
}

func (p *Print) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("print", flag.ExitOnError)

	flagset.StringVar(&p.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&p.split, "split", false, "Create separate file for each track. Defaults to the same directory as the MIDI file.")
	flagset.StringVar(&p.templates, "templates", "", "Loads the formatting templates from a file")
	flagset.BoolVar(&p.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&p.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (p *Print) config() string {
	return p.conf
}

func (p *Print) Execute(smf *midi.SMF) {
	eventlog.EventLog.Verbose = p.verbose
	eventlog.EventLog.Debug = p.debug

	op, err := print.NewPrint()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if p.templates != "" {
		f, err := os.Open(p.templates)
		if err == nil {
			err = op.LoadTemplates(f)
			f.Close()
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	}

	if p.split {
		p.separate(op, smf)
	} else {
		p.write(op, smf)
	}
}

func (p *Print) write(op *print.Print, smf *midi.SMF) {
	out := os.Stdout

	if p.out != "" {
		w, err := os.Create(p.out)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()

		out = w
	}

	err := op.Print(smf, "document", out)
	if err != nil {
		fmt.Printf("Error %v extracting MIDI information\n", err)
	}
}

func (p *Print) separate(op *print.Print, smf *midi.SMF) {
	// Get base filename and Create out directory
	base := strings.TrimSuffix(path.Base(smf.File), path.Ext(smf.File))
	dir := path.Dir(smf.File)

	if p.out != "" {
		dir = p.out
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
	}

	// Print MThd
	filename := fmt.Sprintf("%s.MThd", base)
	w, err := os.Create(path.Join(dir, filename))
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
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
			fmt.Printf("Error: %v", err)
			return
		}

		if err = op.Print(mtrk, "MTrk", w); err != nil {
			fmt.Printf("Error %v extracting MTrk%d information\n", mtrk.TrackNumber, err)
		}

		w.Close()
	}
}
