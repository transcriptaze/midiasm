package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/log"
	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

var cli = []struct {
	cmd     string
	command Command
}{
	{"disassemble", &DISASSEMBLE},
	{"assemble", &ASSEMBLE},
	{"notes", &NOTES},
	{"click", &CLICK},
	{"export", &EXPORT},
	{"transpose", &TRANSPOSE},
	{"help", &HELP},
	{"version", &VERSION},
}

const version = "v0.1.0"

func main() {
	// ... parse command line
	cmd, flagset, err := parse()
	if err != nil {
		fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
		return
	}

	// ... set manufacturer specific options
	if conf := cmd.config(); conf != "" {
		f, err := os.Open(conf)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		manufacturers, err := lib.LoadManufacturers(f)

		f.Close()

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		lib.AddManufacturers(manufacturers)
	}

	// ... set global stuff

	if cmd.Debug() {
		log.SetLogLevel(log.Debug)
	} else if cmd.Verbose() {
		log.SetLogLevel(log.Info)
	}

	context.SetMiddleC(cmd.MiddleC())

	// ... process
	if err := cmd.Execute(flagset); err != nil {
		fmt.Println()
		fmt.Printf("   *** ERROR: %v\n", err)
		fmt.Println()

		os.Exit(1)
	}
}

func parse() (Command, *flag.FlagSet, error) {
	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.cmd == os.Args[1] {
				cmd := c.command
				flagset := cmd.Flagset()
				if err := flagset.Parse(os.Args[2:]); err != nil {
					return cmd, flagset, err
				} else {
					return cmd, flagset, nil
				}
			}
		}
	}

	cmd := &DISASSEMBLE
	flagset := cmd.Flagset()
	if err := flagset.Parse(os.Args[1:]); err != nil {
		return cmd, flagset, err
	}

	return cmd, flagset, nil
}
