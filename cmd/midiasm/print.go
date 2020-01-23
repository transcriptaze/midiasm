package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"github.com/twystd/midiasm/midi/operations"
	"os"
)

type Print struct {
	out     string
	split   bool
	verbose bool
	debug   bool
}

func (p *Print) flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("print", flag.ExitOnError)

	flagset.StringVar(&p.out, "out", "", "Output file path (or directory for split files)")
	flagset.BoolVar(&p.split, "split", false, "Create separate file for each track. Defaults to the same directory as the MIDI file.")
	flagset.BoolVar(&p.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&p.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (p *Print) Execute(smf *midi.SMF) {
	eventlog.EventLog.Verbose = p.verbose
	eventlog.EventLog.Debug = p.debug

	// if p.split {
	// 	p.separate(smf)
	// } else {
	// 	p.write(smf)
	// }

	p.write(smf)
}

func (p *Print) write(smf *midi.SMF) {
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

	op := operations.Print{}
	err := op.Execute(smf, out)
	if err != nil {
		fmt.Printf("Error %v extracting MIDI information\n", err)
	}
}

// func (p *Print) separate(smf *midi.SMF) {
// 	base := strings.TrimSuffix(path.Base(smf.File), path.Ext(smf.File))
// 	dir := path.Dir(smf.File)
//
// 	if p.out != "" {
// 		dir = p.out
// 		err := os.MkdirAll(dir, os.ModeDir)
// 		if err != nil {
// 			fmt.Printf("Error: %v", err)
// 			return
// 		}
// 	}
//
// 	state := struct {
// 		chunk midi.Chunk
// 		w     io.Writer
// 		files []*os.File
// 	}{
// 		files: make([]*os.File, 0),
// 	}
//
// 	defer func() {
// 		for _, f := range state.files {
// 			f.Close()
// 		}
// 	}()
//
// 	f := func(chunk midi.Chunk) (io.Writer, error) {
// 		if chunk != state.chunk {
// 			if _, ok := chunk.(*midi.MThd); ok {
// 				filename := fmt.Sprintf("%s.MThd", base)
// 				if w, err := os.Create(path.Join(dir, filename)); err != nil {
// 					return nil, err
// 				} else {
// 					state.w = w
// 				}
// 			}
//
// 			if mtrk, ok := chunk.(*midi.MTrk); ok {
// 				filename := fmt.Sprintf("%s-%d.MTrk", base, mtrk.TrackNumber)
// 				if w, err := os.Create(path.Join(dir, filename)); err != nil {
// 					return nil, err
// 				} else {
// 					state.w = w
// 				}
// 			}
//
// 			state.chunk = chunk
// 		}
//
// 		return state.w, nil
// 	}
//
// 	q := operations.Print{f}
// 	err := q.Execute(smf)
// 	if err != nil {
// 		fmt.Printf("Error %v extracting notes\n", err)
// 	}
// }
