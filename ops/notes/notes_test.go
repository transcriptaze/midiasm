package notes

import (
	"reflect"
	"testing"
	"time"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

var smf = midi.SMF{
	MThd: &midi.MThd{
		Format:        1,
		Tracks:        2,
		PPQN:          480,
		Division:      480,
		SMPTETimeCode: false,
		SubFrames:     0,
		DropFrame:     false,
	},

	Tracks: []*midi.MTrk{
		&midi.MTrk{
			Events: []*events.Event{
				&events.Event{
					Event: metaevent.MakeTempo(0, 0, 500000),
				},
			},
		},
		&midi.MTrk{
			Events: []*events.Event{
				&events.Event{
					Event: midievent.MakeNoteOn(0, 0, 0, midievent.Note{
						Value: 48,
						Name:  "C3",
						Alias: "C3",
					}, 72),
				},
				&events.Event{
					Event: midievent.MakeNoteOff(480, 480, 0, midievent.Note{
						Value: 48,
						Name:  "C3",
						Alias: "C3",
					}, 64),
				},
			},
		},
	},
}

func TestExtractNotes(t *testing.T) {
	expected := []*Note{
		&Note{
			Channel:       0,
			Note:          48,
			FormattedNote: "C3",
			Velocity:      72,
			StartTick:     0,
			EndTick:       480,
			Start:         0 * time.Millisecond,
			End:           500 * time.Millisecond,
		},
	}

	notes, err := extract(&smf, 0)

	if err != nil {
		t.Fatalf("Error extracting notes from SMF (%v)", err)
	}

	if len(notes) != len(expected) {
		t.Errorf("Incorrectly extracted notes\n   expected:%v\n   got:     %v", expected, notes)
	} else {
		for i := range expected {
			p := expected[i]
			q := notes[i]
			if !reflect.DeepEqual(*p, *q) {
				t.Errorf("Incorrectly extracted note %v\n   expected:%v\n   got:     %v", i+1, *p, *q)
			}
		}
	}

}
