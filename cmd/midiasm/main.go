package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"github.com/twystd/midiasm/midi/eventlog"
	"io/ioutil"
	"os"
)

var cli = map[string]*flag.FlagSet{
	"notes": flag.NewFlagSet("notes", flag.ExitOnError),
}

func main() {
	cmd, filename, err := parse(cli)
	if err != nil {
		fmt.Printf("Unable to parse command line: %v", err)
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		eventlog.Error(fmt.Sprintf("%v", err))
		return
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		eventlog.Error(fmt.Sprintf("%v", err))
		return
	}

	var smf midi.SMF

	if err = smf.UnmarshalBinary(bytes); err != nil {
		eventlog.Error(fmt.Sprintf("%v", err))
		return
	}

	switch cmd {
	case "notes":
		smf.Notes()
	default:
		smf.Render()
	}
}

func parse(cli map[string]*flag.FlagSet) (string, string, error) {
	if len(os.Args) > 1 {
		cmd := os.Args[1]

		if flagset, ok := cli[cmd]; ok {
			if err := flagset.Parse(os.Args[2:]); err != nil {
				return cmd, "", err
			}

			return cmd, flagset.Arg(0), nil
		}
	}

	flag.Parse()

	return "render", flag.Arg(0), nil
}
