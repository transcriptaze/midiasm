package midi

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/events/midi"
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

var programBankMSB = events.Event{
	Bytes: types.Hex{0x00, 0xb7, 0x00, 0x05},
	Event: &midievent.Controller{
		Tag:        "Controller",
		Status:     0xb7,
		Channel:    types.Channel(7),
		Controller: types.Controller{0, "Bank Select (MSB)"},
		Value:      0x05,
	},
}

var programBankLSB = events.Event{
	Bytes: types.Hex{0x00, 0xb7, 0x20, 0x21},
	Event: &midievent.Controller{
		Tag:        "Controller",
		Status:     0xb7,
		Channel:    types.Channel(7),
		Controller: types.Controller{32, "Bank Select (LSB)"},
		Value:      0x21,
	},
}

var programBankMSBCh3 = events.Event{
	Bytes: types.Hex{0x00, 0xb3, 0x00, 0x05},
	Event: &midievent.Controller{
		Tag:        "Controller",
		Status:     0xb3,
		Channel:    types.Channel(3),
		Controller: types.Controller{0, "Bank Select (MSB)"},
		Value:      0x05,
	},
}

var programBankLSBCh5 = events.Event{
	Bytes: types.Hex{0x00, 0xb5, 0x20, 0x21},
	Event: &midievent.Controller{
		Tag:        "Controller",
		Status:     0xb5,
		Channel:    types.Channel(5),
		Controller: types.Controller{32, "Bank Select (LSB)"},
		Value:      0x21,
	},
}

var noteOn = events.Event{
	Bytes: types.Hex{0x00, 0x30, 0x40},
	Event: &midievent.NoteOn{
		Tag:     "NoteOn",
		Status:  0x97,
		Channel: types.Channel(7),
		Note: types.Note{
			Value: 48,
			Name:  "C2",
			Alias: "C2",
		},
		Velocity: 64,
	},
}

var noteOff = events.Event{
	Bytes: types.Hex{0x00, 0x30, 0x40},
	Event: &midievent.NoteOff{
		Tag:     "NoteO.n",
		Status:  0x87,
		Channel: types.Channel(7),
		Note: types.Note{
			Value: 48,
			Name:  "C2",
			Alias: "C2",
		},
		Velocity: 64,
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

func TestValidateProgramBank(t *testing.T) {
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
					// normal program bank
					&programBankMSB,
					&programBankLSB,
					&noteOn,
					&noteOff,

					// missing program bank LSB
					&programBankMSB,
					&noteOn,

					// missing program bank MSB
					&programBankLSB,
					&noteOff,

					// unmatched channel
					&programBankMSBCh3,
					&programBankLSBCh5,
					&noteOn,
					&noteOff,

					&endOfTrack,
				},
			},
		},
	}

	expected := []ValidationError{
		ValidationError(fmt.Errorf("Track 1: 'Bank Select MSB' event @5 missing LSB (Controller)")),
		ValidationError(fmt.Errorf("Track 1: 'Bank Select LSB' event @7 missing MSB (Controller)")),
		ValidationError(fmt.Errorf("Track 1: 'Bank Select MSB' event @9 LSB on another channel (Controller)")),
		ValidationError(fmt.Errorf("Track 1: 'Bank Select LSB' event @10 MSB on another channel (Controller)")),
	}

	diff(t, expected, smf.Validate())
}

func diff(t *testing.T, expected, errors []ValidationError) {
	if len(errors) != len(expected) {
		t.Errorf("Validation returned %d errors, expected: %v", len(errors), len(expected))
	}

loop1:
	for _, e := range expected {
		for _, err := range errors {
			if reflect.DeepEqual(err, e) {
				continue loop1
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
