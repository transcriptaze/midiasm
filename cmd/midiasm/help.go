package main

import (
	"flag"
)

type Help struct {
	command
}

var HELP = Help{}

func (t *Help) flagset() *flag.FlagSet {
	return flag.NewFlagSet("help", flag.ExitOnError)
}

func (t Help) Execute(filename string) error {
	return nil
}
