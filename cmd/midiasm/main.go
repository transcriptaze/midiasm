package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"io/ioutil"
	"os"
)

var notes bool

func main() {
	flag.BoolVar(&notes, "notes", false, "Extract notes from MIDI sequence")
	flag.Parse()

	filename := flag.Arg(0)

	f, err := os.Open(filename)
	if err != nil {
		eventlog.Error(fmt.Sprintf("%v", err))
		return
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		eventlog.Error(fmt.Sprintf("%v", err))
		return
	}

	var smf midi.SMF

	if err = smf.UnmarshalBinary(bytes); err != nil {
		eventlog.Error(fmt.Sprintf("%v", err))
		return
	}

	if notes {
		smf.Notes()
	} else {
		smf.Render()
	}
}
