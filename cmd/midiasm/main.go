package main

import (
	"fmt"
	"github.com/twystd/midiasm/midi"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println()
	fmt.Println("MIDIASM v0.00.0")
	fmt.Println()

	filename := os.Args[1]

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

	smf.Render()
}
