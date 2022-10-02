package main

import (
	"flag"
	"fmt"
)

type Help struct {
	command
}

var HELP = Help{}

func (t *Help) flagset() *flag.FlagSet {
	return flag.NewFlagSet("help", flag.ExitOnError)
}

func (t Help) Execute() error {
	fmt.Println()
	fmt.Println("  Usage: midiasm <command> <options>")
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range cli {
		fmt.Printf("    %v\n", c.cmd)
	}

	fmt.Println()
	fmt.Println("  Defaults to 'disassemble' if the command is not provided.")
	fmt.Println()
	fmt.Println("  Use 'midiasm help <command>' for command specific information.")
	fmt.Println()

	return nil
}
