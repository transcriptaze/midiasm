package processors

import (
	"fmt"
	"github.com/twystd/midiasm/midi"
	"io"
)

type Print struct {
	Writer io.Writer
}

func (x *Print) Execute(smf *midi.SMF) error {
	smf.Header.Print(x.Writer)

	fmt.Fprintln(x.Writer)
	fmt.Fprintln(x.Writer)

	for _, track := range smf.Tracks {
		track.Print(x.Writer)
	}

	return nil
}
