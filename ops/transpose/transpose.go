package transpose

import (
	"bytes"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

type Transpose struct {
}

func (t *Transpose) Execute(smf *midi.SMF, steps int) ([]byte, error) {
	for i, mtrk := range smf.Tracks {
		transposed := transpose(*mtrk, steps)
		smf.Tracks[i] = transposed
	}

	var b bytes.Buffer
	var e = midifile.NewEncoder(&b)

	if err := e.Encode(*smf); err != nil {
		return nil, err
	} else {
		return b.Bytes(), nil
	}
}

func transpose(mtrk midi.MTrk, steps int) *midi.MTrk {
	track := midi.MTrk{
		Tag:         "MTrk",
		TrackNumber: mtrk.TrackNumber,
		Events:      []*events.Event{},
	}

	for i, _ := range mtrk.Events {
		event := mtrk.Events[i]
		switch v := event.Event.(type) {
		case *metaevent.KeySignature:
		case metaevent.KeySignature:
			track.Events = append(track.Events, &events.Event{
				Event: v.Transpose(mtrk.Context, steps),
			})

		case *midievent.NoteOn:
		case midievent.NoteOn:
			track.Events = append(track.Events, &events.Event{
				Event: v.Transpose(mtrk.Context, steps),
			})

		case *midievent.NoteOff:
		case midievent.NoteOff:
			track.Events = append(track.Events, &events.Event{
				Event: v.Transpose(mtrk.Context, steps),
			})

		default:
			track.Events = append(track.Events, event)
		}
	}

	return &track
}
