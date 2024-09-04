package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
)

type Command interface {
	Flagset(flagset *flag.FlagSet) *flag.FlagSet
	Execute(flagset *flag.FlagSet) error
	Help()
}

func decode(filename string) (*midi.SMF, error) {
	var r io.Reader

	if b, err := read(filename); err != nil {
		return nil, err
	} else {
		r = bytes.NewReader(b)
	}

	decoder := midifile.NewDecoder()

	if smf, err := decoder.Decode(r); err != nil {
		return nil, err
	} else if smf == nil {
		return nil, fmt.Errorf("failed to decode MIDI file")
	} else {
		return smf, nil
	}
}

func read(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}
