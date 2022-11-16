package transpose

import (
	"bytes"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/encoding/midi"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

type Transpose struct {
}

func (t *Transpose) Execute(smf *midi.SMF, steps int) ([]byte, error) {
	for _, mtrk := range smf.Tracks {
		transpose(mtrk, steps)
	}

	var b bytes.Buffer
	var e = midifile.NewEncoder(&b)

	if err := e.Encode(*smf); err != nil {
		return nil, err
	} else {
		return b.Bytes(), nil
	}
}

func transpose(mtrk *midi.MTrk, steps int) {
	for _, event := range mtrk.Events {
		switch v := event.Event.(type) {
		case *metaevent.KeySignature:
			v.Transpose(mtrk.Context, steps)

		case *midievent.NoteOn:
			v.Transpose(mtrk.Context, steps)

		case *midievent.NoteOff:
			v.Transpose(mtrk.Context, steps)
		}
	}
}
