package main

import (
	"flag"
	"fmt"
	"github.com/twystd/midiasm/midi"
	"io/ioutil"
	"os"
)

var notes bool

func main() {
	fmt.Println()
	fmt.Println("MIDIASM v0.00.0")
	fmt.Println()

	flag.BoolVar(&notes, "notes", false, "Extract notes from MIDI sequence")
	flag.Parse()

	filename := flag.Arg(0)

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println()
		fmt.Printf("   Error opening %s: %v", filename, err)
		return
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("   Error reading %s: %v", filename, err)
		fmt.Println()
		return
	}

	var smf midi.SMF

	if err = smf.UnmarshalBinary(bytes); err != nil {
		fmt.Printf("   Error reading %s: %v", filename, err)
		fmt.Println()
		return
	}

	if notes {
		smf.Notes()
	} else {
		smf.Render()
	}
}
