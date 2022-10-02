package main

import (
	"flag"
	"fmt"
)

type Version struct {
	command
}

var VERSION = Version{}

func (t *Version) flagset() *flag.FlagSet {
	return flag.NewFlagSet("version", flag.ExitOnError)
}

func (t Version) Execute() error {
	fmt.Printf("midiasm %v\n", version)
	return nil
}
