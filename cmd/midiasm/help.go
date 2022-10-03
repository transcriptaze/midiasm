package main

import (
	"flag"
	"fmt"
)

type Help struct {
	command
}

var HELP = Help{}

func (h *Help) Flagset() *flag.FlagSet {
	flagset := flag.NewFlagSet("help", flag.ExitOnError)

	h.flags = flagset

	return flagset
}

func (h Help) Help() {
}

func (h Help) Execute() error {
	for _, c := range cli {
		if c.cmd == h.flags.Arg(0) {
			c.command.Help()
			return nil
		}
	}

	h.help()

	return nil
}

func (h Help) help() {
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
}
