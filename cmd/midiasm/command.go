package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/lib"
)

	Execute(flagset *flag.FlagSet) error
type command struct {
	conf    string
	c4      bool
	verbose bool
	debug   bool
}

func (c command) Config() string {
	return c.conf
}

func (c command) MiddleC() lib.MiddleC {
	if c.c4 {
		return lib.C4
	}

	return lib.C3
}

func (c command) Debug() bool {
	return c.debug
}

func (c command) Verbose() bool {
	return c.verbose
}

func (c *command) flagset(name string) *flag.FlagSet {
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	flagset.BoolVar(&c.c4, "C4", c.c4, "Sets middle C to C4 (Yamaho convention). Defaults to C3")
	flagset.BoolVar(&c.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&c.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (c *command) decode(filename string) (*midi.SMF, error) {
	var r io.Reader

	if b, err := c.read(filename); err != nil {
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

func (c *command) read(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}
