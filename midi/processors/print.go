package processors

import (
	"fmt"
	"github.com/twystd/midiasm/midi"
	"io"
)

type Print struct {
	Writer io.Writer
}

func (p *Print) Execute(smf *midi.SMF) {
	smf.Header.Print(p.Writer)

	fmt.Fprintln(p.Writer)
	fmt.Fprintln(p.Writer)

	for _, track := range smf.Tracks {
		track.Print(p.Writer)
	}
}
