package main

import (
	"flag"
	"fmt"
)

type version struct {
}

var Version = version{}

func (v *version) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (v version) Help() {
}

func (v version) Execute(flagset *flag.FlagSet) error {
	fmt.Printf("midiasm %v\n", VERSION)

	return nil
}
