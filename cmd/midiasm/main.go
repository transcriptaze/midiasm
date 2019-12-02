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
	"notes": notesFlagset(),
}

var options = struct {
	out     string
	verbose bool
	debug   bool
}{}

func main() {
	cmd, filename, err := parse(cli)
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

	switch cmd {
	case "notes":
		notes(&smf)
	default:
		render(&smf)
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

	flag.StringVar(&options.out, "out", "", "Output file path")
	flag.BoolVar(&options.verbose, "verbose", false, "Enable progress information")
	flag.BoolVar(&options.debug, "debug", false, "Enable debugging information")
	flag.Parse()

	return "render", flag.Arg(0), nil
}

func render(smf *midi.SMF) {
	w := os.Stdout
	err := error(nil)

	if options.out != "" {
		w, err = os.Create(options.out)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		defer w.Close()
	}

	eventlog.EventLog.Verbose = options.verbose
	eventlog.EventLog.Debug = options.debug

	smf.Render(w)
}
