package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"io/ioutil"
	"os"
)

type command interface {
	flagset() *flag.FlagSet
	Execute(*midi.SMF)
}

var cli = map[string]command{
	"print": &Print{},
	"notes": &Notes{},
}

var options = struct {
	out     string
	verbose bool
	debug   bool
}{}

func main() {
	cmd, filename, err := parse()
	if err != nil {
		fmt.Printf("Error: unable to parse command line (%v)\n", err)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var smf midi.SMF

	if err = smf.UnmarshalBinary(bytes); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	cmd.Execute(&smf)
}

func parse() (command, string, error) {
	cmd := &Print{}
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
	if err := flagset.Parse(os.Args[2:]); err != nil {
		return cmd, "", err
	}

	return cmd, flagset.Arg(0), nil
}
