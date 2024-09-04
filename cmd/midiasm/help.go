package main

import (
	"flag"
	"fmt"
)

type help struct {
}

var Help = help{}

func (h *help) Flagset(flagset *flag.FlagSet) *flag.FlagSet {
	return flagset
}

func (h help) Help() {
}

func (h help) Execute(flagset *flag.FlagSet) error {
	for _, c := range cli {
		if c.cmd == flagset.Arg(0) {
			c.command.Help()
			return nil
		}
	}

	h.help()

	return nil
}

func (h help) help() {
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
