package main

import (
	"flag"
	"fmt"
)

type Version struct {
	command
}

var VERSION = Version{}

func (v *Version) Flagset() *flag.FlagSet {
	return flag.NewFlagSet("version", flag.ExitOnError)
}

func (v Version) Help() {
}

func (v Version) Execute(flagset *flag.FlagSet) error {
	fmt.Printf("midiasm %v\n", version)

	return nil
}
