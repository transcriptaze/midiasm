package midi

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/types"
	"reflect"
	"testing"
)

var tempo = events.Event{
	Bytes: types.Hex{0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20},
	Event: &metaevent.Tempo{
		Tag:    "Tempo",
		Status: 0xff,
		Type:   types.MetaEventType(0x51),
		Tempo:  500000,
	},
}

var smpteOffset = events.Event{
	Bytes: types.Hex{0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27},
	Event: &metaevent.SMPTEOffset{
		Tag:              "SMPTEOffset",
		Status:           0xff,
		Type:             types.MetaEventType(0x54),
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	},
}

var endOfTrack = events.Event{
	Bytes: types.Hex{0x00, 0xff, 0x2f, 0x00},
	Event: &metaevent.EndOfTrack{
		Tag:    "EndOfTrack",
		Status: 0xff,
		Type:   types.MetaEventType(0x2f),
	},
}

func TestValidateFormat1(t *testing.T) {
	smf := SMF{
		MThd: &MThd{
			Length: 6,
			Format: 1,
		},

		Tracks: []*MTrk{
			&MTrk{
				TrackNumber: 0,
				Events: []*events.Event{
					&endOfTrack,
				},
			},

			&MTrk{
				TrackNumber: 1,
				Events: []*events.Event{
					&tempo,
					&smpteOffset,
					&endOfTrack,
				},
			},
		},
	}

	expected := []ValidationError{
		ValidationError(fmt.Errorf("Track 1: unexpected event (Tempo)")),
		ValidationError(fmt.Errorf("Track 1: unexpected event (SMPTEOffset)")),
	}

	errors := smf.Validate()
	if len(errors) != len(expected) {
		t.Errorf("Validation returned %d errors, expected: %v", len(errors), len(expected))
	}

loop:
	for _, e := range expected {
		for _, err := range errors {
			if reflect.DeepEqual(err, e) {
				continue loop
			}
		}
		t.Errorf("Missing expected error: %v", e)
	}

loop2:
	for _, e := range errors {
		for _, err := range expected {
			if reflect.DeepEqual(err, e) {
				continue loop2
			}
		}
		t.Errorf("Unexpected error: %v", e)
	}
}
