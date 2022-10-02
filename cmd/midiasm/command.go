package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/encoding/midifile"
	"github.com/transcriptaze/midiasm/midi/types"
)

type Command interface {
	flagset() *flag.FlagSet
	config() string
	MiddleC() types.MiddleC
	Execute() error
}

type command struct {
	conf    string
	c4      bool
	verbose bool
	debug   bool
}

func (c command) config() string {
	return c.conf
}

func (c command) MiddleC() types.MiddleC {
	if c.c4 {
		return types.C4
	}

	return types.C3
}

func (c *command) flagset(name string) *flag.FlagSet {
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	flagset.BoolVar(&c.c4, "C4", c.c4, "Sets middle C to C4 (Yamaho convention). Defaults to C3")
	flagset.BoolVar(&c.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&c.debug, "debug", false, "Enable debugging information")

	return flagset
}

func (c *command) decode(filename string) (*midi.SMF, error) {
	bytes, err := c.read(filename)
	if err != nil {
		return nil, err
	}

	decoder := midifile.NewDecoder()

	smf, err := decoder.Decode(bytes)
	if err != nil {
		return nil, err
	}

	if smf == nil {
		return nil, fmt.Errorf("failed to decode MIDI file")
	}

	smf.File = filename

	return smf, nil
}

func (c *command) read(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}
