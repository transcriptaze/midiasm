package main

import (
	"flag"
)

type Version struct {
	command
}

var VERSION = Version{}

func (t *Version) flagset() *flag.FlagSet {
	return flag.NewFlagSet("version", flag.ExitOnError)
}

func (t Version) Execute(filename string) error {
	return nil
}
