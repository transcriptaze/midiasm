package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"io/ioutil"
	"os"
)

var cli = map[string]*flag.FlagSet{
	"notes": flag.NewFlagSet("notes", flag.ExitOnError),
}

var options = struct {
	out string
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
		smf.Notes()
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
	flag.StringVar(&options.out, "o", "", "Output file path")
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

	smf.Render(w)
}
