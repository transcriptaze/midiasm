package processors

import (
	"github.com/twystd/midiasm/midi"
	"io"
)

type Print struct {
	Writer func(midi.Chunk) (io.Writer, error)
}

func (p *Print) Execute(smf *midi.SMF) error {
	if w, err := p.Writer(smf.Header); err != nil {
		return err
	} else {
		smf.Header.Print(w)
	}

	for _, track := range smf.Tracks {
		if w, err := p.Writer(track); err != nil {
			return err
		} else {
			track.Print(w)
		}
	}

	return nil
}
