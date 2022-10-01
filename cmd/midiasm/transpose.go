package main

import (
	"flag"

	"github.com/transcriptaze/midiasm/midi/types"
	"github.com/transcriptaze/midiasm/ops/transpose"
)

type Transpose struct {
	command
	out       string
	transpose int
}

var TRANSPOSE = Transpose{
	command: command{
		middleC: types.C3,
	},
}

func (t *Transpose) flagset() *flag.FlagSet {
	flagset := t.command.flagset("transpose")

	flagset.StringVar(&t.out, "out", "", "Output file path")
	flagset.IntVar(&t.transpose, "transpose", 0, "Transpose notes up or down")

	return flagset
}

func (t Transpose) Execute(filename string) error {
	return t.execute()
}

func (t Transpose) execute() error {
	// var w = os.Stdout
	// var err error

	// if t.out != "" {
	// 	if w, err = os.Create(t.out); err != nil {
	// 		fmt.Printf("Error: %v\n", err)
	// 		return
	// 	}

	// 	defer w.Close()
	// }

	// eventlog.EventLog.Verbose = t.verbose
	// eventlog.EventLog.Debug = t.debug

	// p := transpose.Transpose{w}
	// if err = p.Execute(smf); err != nil {
	// 	fmt.Printf("Error transposing %q (%v)\n", smf.File, err)
	// }

	op := transpose.Transpose{}

	return op.Execute()
}
