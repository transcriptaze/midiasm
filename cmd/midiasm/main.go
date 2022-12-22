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

var options = struct {
	conf    string
	c4      bool
	verbose bool
	debug   bool
}{}

const version = "v0.1.0"

func main() {
	// ... parse command line
	cmd, flagset, err := parse()
	if err != nil {
		fmt.Printf("ERROR: unable to parse command line (%v)\n", err)
		return
	}

	// ... set manufacturer specific options
	if options.conf != "" {
		f, err := os.Open(options.conf)
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
	if options.debug {
		log.SetLogLevel(log.Debug)
	} else if options.verbose {
		log.SetLogLevel(log.Info)
	}

	if options.c4 {
		context.SetMiddleC(lib.C4)
	} else {
		context.SetMiddleC(lib.C3)
	}

	// ... process
	if err := cmd.Execute(flagset); err != nil {
		fmt.Println()
		fmt.Printf("   *** ERROR: %v\n", err)
		fmt.Println()

		os.Exit(1)
	}
}

func parse() (Command, *flag.FlagSet, error) {
	flagset := flag.NewFlagSet("midiasm", flag.ExitOnError)

	flagset.BoolVar(&options.c4, "C4", options.c4, "Sets middle C to C4 (Yamaho convention). Defaults to C3")
	flagset.BoolVar(&options.verbose, "verbose", false, "Enable progress information")
	flagset.BoolVar(&options.debug, "debug", false, "Enable debugging information")

	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.cmd == os.Args[1] {
				cmd := c.command
				flagset = cmd.Flagset(flagset)
				if err := flagset.Parse(os.Args[2:]); err != nil {
					return cmd, flagset, err
				} else {
					return cmd, flagset, nil
				}
			}
		}
	}

	cmd := &DISASSEMBLE
	flagset = cmd.Flagset(flagset)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		return cmd, flagset, err
	}

	return cmd, flagset, nil
}
