package main

import (
	"fmt"
	"os"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

var cli = map[string]Command{
	"disassemble": &DISASSEMBLE,
	"notes":       &NOTES,
	"click":       &CLICK,
	"export":      &EXPORT,
	"transpose":   &TRANSPOSE,
	"help":        &HELP,
	"version":     &VERSION,
}

func main() {
	// ... parse command line
	cmd, filename, err := parse()
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

	// ... set middle C convention
	context.SetMiddleC(cmd.MiddleC())

	// ... process
	if err := cmd.Execute(filename); err != nil {
		fmt.Println()
		fmt.Printf("   *** ERROR: %v\n", err)
		fmt.Println()
	}
}

func parse() (Command, string, error) {
	cmd := &DISASSEMBLE
	if len(os.Args) > 1 {
		c, ok := cli[os.Args[1]]
		if ok {
			flagset := c.flagset()
			if err := flagset.Parse(os.Args[2:]); err != nil {
				return cmd, "", err
			}

			return c, flagset.Arg(0), nil
		}
	}

	flagset := cmd.flagset()
	if err := flagset.Parse(os.Args[1:]); err != nil {
		return cmd, "", err
	}

	return cmd, flagset.Arg(0), nil
}
