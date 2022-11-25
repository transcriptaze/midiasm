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

var trackname = makeEvent(
	metaevent.MakeTrackName(0, 0, "Example 1"),
	[]byte{0x0, 0xff, 0x3, 0x9, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31}...)

var keysignatureFSM = makeEvent(
	metaevent.MakeKeySignature(0, 0, 6, types.Major, "F♯ major", []byte{0x00, 0xff, 0x59, 0x02, 0x06, 0x00}...),
	[]byte{0x00, 0xff, 0x59, 0x02, 0x06, 0x00}...)

var keysignatureEFm = makeEvent(
	metaevent.MakeKeySignature(0, 0, -6, types.Minor, "E♭ minor", []byte{0x00, 0xff, 0x59, 0x02, 0xfa, 0x01}...),
	[]byte{0x00, 0xff, 0x59, 0x02, 0xfa, 0x01}...)

var noteOnC3v72 = makeEvent(
	midievent.MakeNoteOn(0, 0, 1, midievent.Note{Value: 48, Name: "C3", Alias: "C3"}, 72, []byte{0x00, 0x91, 0x30, 0x48}...),
	[]byte{0x00, 0x91, 0x30, 0x48}...)

var noteOnC3v0 = makeEvent(
	midievent.MakeNoteOn(0, 0, 1, midievent.Note{Value: 48, Name: "C3", Alias: "C3"}, 0, []byte{0x00, 0x30, 0x00}...),
	[]byte{0x00, 0x30, 0x00}...)

var noteOnC3v64 = makeEvent(
	midievent.MakeNoteOn(0, 0, 1, midievent.Note{Value: 48, Name: "C3", Alias: "C3"}, 64, []byte{0x00, 0x30, 0x40}...),
	[]byte{0x00, 0x30, 0x40}...)

var noteOnC3v32 = makeEvent(
	midievent.MakeNoteOn(0, 0, 1, midievent.Note{Value: 48, Name: "C3", Alias: "C3"}, 32, []byte{0x00, 0x30, 0x20}...),
	[]byte{0x00, 0x30, 0x20}...)

var noteOnCS3 = makeEvent(
	midievent.MakeNoteOn(0, 0, 1, midievent.Note{Value: 49, Name: "C♯3", Alias: "C♯3"}, 72, []byte{0x00, 0x91, 0x31, 0x48}...),
	[]byte{0x00, 0x91, 0x31, 0x48}...)

var noteOffCS3Alias = makeMidiEvent(
	midievent.MakeNoteOff(0, 0, 1, midievent.Note{Value: 49, Name: "C♯3", Alias: "D♭3"}, 100, []byte{0x00, 0x81, 0x31, 0x64}...),
	[]byte{0x00, 0x81, 0x31, 0x64}...)

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
			t.Errorf("\n   expected:%#v\n   got:     %#v", e.Event, mtrk.Events[i].Event)
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
