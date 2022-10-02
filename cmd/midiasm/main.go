package main

import (
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/log"
	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

var cli = []struct {
	cmd     string
	command Command
}{
	{"disassemble", &DISASSEMBLE},
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
	cmd, err := parse()
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

		manufacturers, err := types.LoadManufacturers(f)

		f.Close()

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		types.AddManufacturers(manufacturers)
	}

	// ... set global stuff

	if cmd.Debug() {
		log.SetLogLevel(log.Debug)
	} else if cmd.Verbose() {
		log.SetLogLevel(log.Info)
	}

	context.SetMiddleC(cmd.MiddleC())

	// ... process
	if err := cmd.Execute(); err != nil {
		fmt.Println()
		fmt.Printf("   *** ERROR: %v\n", err)
		fmt.Println()

		os.Exit(1)
	}
}

func parse() (Command, error) {
	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.cmd == os.Args[1] {
				cmd := c.command
				flagset := cmd.flagset()
				if err := flagset.Parse(os.Args[2:]); err != nil {
					return cmd, err
				} else {
					return cmd, nil
				}
			}
		}
	}

	cmd := &DISASSEMBLE
	flagset := cmd.flagset()
	if err := flagset.Parse(os.Args[1:]); err != nil {
		return cmd, err
	}

	return cmd, nil
}
