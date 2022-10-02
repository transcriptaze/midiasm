package midi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/types"
)

var keysignatureFSM = events.Event{
	Bytes: types.Hex{0x00, 0xff, 0x59, 0x02, 0x06, 0x00},
	Event: &metaevent.KeySignature{
		Tag:         "KeySignature",
		Status:      0xff,
		Type:        types.MetaEventType(0x59),
		Accidentals: 6,
		KeyType:     0,
		Key:         "F♯ major",
	},
}

var keysignatureEFm = events.Event{
	Bytes: types.Hex{0x00, 0xff, 0x59, 0x02, 0xfa, 0x01},
	Event: &metaevent.KeySignature{
		Tag:         "KeySignature",
		Status:      0xff,
		Type:        types.MetaEventType(0x59),
		Accidentals: -6,
		KeyType:     1,
		Key:         "E♭ minor",
	},
}

var noteOnC3v72 = events.Event{
	Bytes: types.Hex{0x00, 0x91, 0x30, 0x48},
	Event: &midievent.NoteOn{
		Tag:     "NoteOn",
		Status:  0x91,
		Channel: types.Channel(0x01),
		Note: midievent.Note{
			Value: 48,
			Name:  "C3",
			Alias: "C3",
		},
		Velocity: 72,
	},
}

var noteOnC3v0 = events.Event{
	Bytes: types.Hex{0x00, 0x30, 0x00},
	Event: &midievent.NoteOn{
		Tag:     "NoteOn",
		Status:  0x91,
		Channel: types.Channel(0x01),
		Note: midievent.Note{
			Value: 48,
			Name:  "C3",
			Alias: "C3",
		},
		Velocity: 0,
	},
}

var noteOnC3v64 = events.Event{
	Bytes: types.Hex{0x00, 0x30, 0x40},
	Event: &midievent.NoteOn{
		Tag:     "NoteOn",
		Status:  0x91,
		Channel: types.Channel(0x01),
		Note: midievent.Note{
			Value: 48,
			Name:  "C3",
			Alias: "C3",
		},
		Velocity: 64,
	},
}

var noteOnC3v32 = events.Event{
	Bytes: types.Hex{0x00, 0x30, 0x20},
	Event: &midievent.NoteOn{
		Tag:     "NoteOn",
		Status:  0x91,
		Channel: types.Channel(0x01),
		Note: midievent.Note{
			Value: 48,
			Name:  "C3",
			Alias: "C3",
		},
		Velocity: 32,
	},
}

var noteOnCS3 = events.Event{
	Bytes: types.Hex{0x00, 0x91, 0x31, 0x48},
	Event: &midievent.NoteOn{
		Tag:     "NoteOn",
		Status:  0x91,
		Channel: types.Channel(0x01),
		Note: midievent.Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		},
		Velocity: 72,
	},
}

var noteOffCS3Alias = events.Event{
	Bytes: types.Hex{0x00, 0x81, 0x31, 0x64},
	Event: &midievent.NoteOff{
		Tag:     "NoteOff",
		Status:  0x81,
		Channel: types.Channel(0x01),
		Note: midievent.Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "D♭3",
		},
		Velocity: 100,
	},
}

func TestUnmarshalNoteAlias(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x14,
		0x00, 0xff, 0x59, 0x02, 0x06, 0x00,
		0x00, 0x91, 0x31, 0x48,
		0x00, 0xff, 0x59, 0x02, 0xfa, 0x01,
		0x00, 0x81, 0x31, 0x64,
	}

	expected := MTrk{
		Tag:         "MTrk",
		TrackNumber: 1,
		Length:      74,
		Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x4a},
		Events: []*events.Event{
			&keysignatureFSM,
			&noteOnCS3,
			&keysignatureEFm,
			&noteOffCS3Alias,
		},
	}

	mtrk := MTrk{
		TrackNumber: 1,
		Context:     context.NewContext(),
	}

	if err := mtrk.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling MTrk: %v", err)
	}

	for i, e := range expected.Events {
		if !reflect.DeepEqual(e, mtrk.Events[i]) {
			t.Errorf("MTrk: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", e, mtrk.Events[i])
		}
	}
}

func TestUnmarshalWithRunningStatus(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x18,
		0x00, 0x91, 0x30, 0x48,
		0x00, 0x30, 0x00,
		0x00, 0x30, 0x40,
		0x00, 0x30, 0x20,
		0x00, 0xff, 0x2f, 0x00,
	}

	expected := []*events.Event{
		&noteOnC3v72,
		&noteOnC3v0,
		&noteOnC3v64,
		&noteOnC3v32,
		&endOfTrack,
	}

	mtrk := MTrk{
		TrackNumber: 1,
		Context:     context.NewContext(),
	}

	if err := mtrk.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling MTrk: %v", err)
	}

	for i, e := range expected {
		if !reflect.DeepEqual(e, mtrk.Events[i]) {
			t.Errorf("MTrk: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", e, mtrk.Events[i])
		}
	}
}

func TestUnmarshalWithInvalidRunningStatus(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0b,
		0x00, 0x3c, 0x4c,
		0x00, 0x91, 0x31, 0x48,
		0x00, 0xff, 0x2f, 0x00,
	}

	expected := fmt.Errorf("Unrecognised MIDI event: 30")

	mtrk := MTrk{
		TrackNumber: 1,
		Context:     context.NewContext(),
	}

	err := mtrk.UnmarshalBinary(bytes)
	if err == nil {
		t.Fatalf("Expected error unmarshaling MTrk - got: %v", nil)
	}

	if !reflect.DeepEqual(err, expected) {
		t.Fatalf("Incorrect error unmarshaling SMF:\nexpected: %+v\n     got: %+v", expected, err)
	}
}
