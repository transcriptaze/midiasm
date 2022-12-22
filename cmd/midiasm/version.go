package main

import (
	"flag"
	"fmt"
)

type Version struct {
}

var VERSION = Version{}

func (v *Version) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (v Version) Help() {
}

func (v Version) Execute(flagset *flag.FlagSet) error {
	fmt.Printf("midiasm %v\n", version)

	return nil
}
