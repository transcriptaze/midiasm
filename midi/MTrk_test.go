package midi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/io"
)

var trackname = events.NewEvent(0, 0, nil, []byte{})
var keysignatureFSM = events.NewEvent(0, 0, nil, []byte{0x00, 0xff, 0x59, 0x02, 0x06, 0x00})
var keysignatureEFm = events.NewEvent(0, 0, nil, []byte{0x00, 0xff, 0x59, 0x02, 0xfa, 0x01})
var noteOnC3v72 = events.NewEvent(0, 0, nil, []byte{0x00, 0x91, 0x30, 0x48})
var noteOnC3v0 = events.NewEvent(0, 0, nil, []byte{0x00, 0x30, 0x00})
var noteOnC3v64 = events.NewEvent(0, 0, nil, []byte{0x00, 0x30, 0x40})
var noteOnC3v32 = events.NewEvent(0, 0, nil, []byte{0x00, 0x30, 0x20})
var noteOnCS3 = events.NewEvent(0, 0, nil, []byte{0x00, 0x91, 0x31, 0x48})
var noteOffCS3Alias = events.NewEvent(0, 0, nil, []byte{0x00, 0x81, 0x31, 0x64})

func init() {
	trackname.Event, _ = metaevent.NewTrackName(0, 0, []byte("Example 1"))
	keysignatureFSM.Event, _ = metaevent.NewKeySignature(nil, 0, 0, []byte{0x06, 0x00})
	keysignatureEFm.Event, _ = metaevent.NewKeySignature(nil, 0, 0, []byte{0xfa, 0x01})

	noteOnC3v72.Event, _ = midievent.NewNoteOn(nil, 0, 0, IO.TestReader([]byte{0x00, 0x91}, []byte{0x30, 0x48}), 0x91)
	noteOnC3v0.Event, _ = midievent.NewNoteOn(nil, 0, 0, IO.TestReader([]byte{0x00}, []byte{0x30, 0x00}), 0x91)  // running status
	noteOnC3v64.Event, _ = midievent.NewNoteOn(nil, 0, 0, IO.TestReader([]byte{0x00}, []byte{0x30, 0x40}), 0x91) // running status
	noteOnC3v32.Event, _ = midievent.NewNoteOn(nil, 0, 0, IO.TestReader([]byte{0x00}, []byte{0x30, 0x20}), 0x91) // running status
	noteOnCS3.Event, _ = midievent.NewNoteOn(nil, 0, 0, IO.TestReader([]byte{0x00, 0x91}, []byte{0x31, 0x48}), 0x91)

	noteOffCS3Alias.Event, _ = midievent.NewNoteOff(nil, 0, 0, IO.TestReader([]byte{0x00, 0x81}, []byte{0x31, 0x64}), 0x81)
}

func TestMTrkMarshalTrack0(t *testing.T) {
	expected := []byte{
		0x4d, 0x54, 0x72, 0x6b, // MTrk
		0x00, 0x00, 0x00, 0x11, // length
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31, // TrackName
		//  0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20, 0x00, 0xff, 0x2f, 0x00, 0x4d, 0x54
		//       0x72, 0x6b, 0x00, 0x00, 0x00, 0x58, 0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74
		//       0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72, 0x00, 0xc0, 0x19, 0x00, 0xff, 0x58, 0x04
		//       0x04, 0x02, 0x18, 0x08, 0x00, 0x90, 0x30, 0x48, 0x00, 0xff, 0x59, 0x02, 0x00, 0x01, 0x00, 0xb0
		//       0x65, 0x00, 0x00, 0xb0, 0x64, 0x00, 0x00, 0xb0, 0x06, 0x06, 0x83, 0x60, 0x80, 0x30, 0x40, 0x00
		//       0x90, 0x32, 0x48, 0x83, 0x60, 0x80, 0x32, 0x40, 0x00, 0x90, 0x34, 0x48, 0x83, 0x60, 0x80, 0x34
		//       0x40, 0x00, 0x90, 0x35, 0x48, 0x83, 0x60, 0x80, 0x35, 0x40,
		0x00, 0xff, 0x2f, 0x00, // EndOfTrack
	}

	mtrk := MTrk{
		Tag:    "MTrk",
		Length: 24,
		Events: []*events.Event{
			trackname,
			endOfTrack,
		},
	}

	bytes, err := mtrk.MarshalBinary()
	if err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(bytes, expected) {
		t.Errorf("incorrectly marshalled\n   expected:%#v\n   got:     %#v", expected, bytes)
	}
}

func TestUnmarshalNoteAlias(t *testing.T) {
	t.Skip() // FIXME can't do note off aliases until MTrk parser is reworked
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
			keysignatureFSM,
			noteOnCS3,
			keysignatureEFm,
			noteOffCS3Alias,
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
		noteOnC3v72,
		noteOnC3v0,
		noteOnC3v64,
		noteOnC3v32,
		endOfTrack,
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
