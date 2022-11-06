package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParseNoteOnInMajorKey(t *testing.T) {
	expected := NoteOn{
		event: event{
			tick:    2400,
			delta:   480,
			bytes:   []byte{0x00, 0x91, 0x31, 0x48},
			tag:     types.TagNoteOn,
			Status:  0x91,
			Channel: 1,
		},
		Note: Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		},
		Velocity: 72,
	}

	ctx := context.NewContext()
	r := IO.TestReader([]byte{0x00, 0x91}, []byte{0x31, 0x48})

	event, err := Parse(2400, 480, r, 0x91, ctx)
	if err != nil {
		t.Fatalf("Unexpected NoteOn event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOn event parse error - returned %v", event)
	}

	event, ok := event.(*NoteOn)
	if !ok {
		t.Fatalf("NoteOn event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid NoteOn event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestParseNoteOnInMinorKey(t *testing.T) {
	expected := NoteOn{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0x91, 0x31, 0x48},

			tag:     types.TagNoteOn,
			Status:  0x91,
			Channel: 1,
		},
		Note: Note{
			Value: 49,
			Name:  "D♭3",
			Alias: "D♭3",
		},
		Velocity: 72,
	}

	ctx := context.NewContext().UseFlats()
	r := IO.TestReader([]byte{0x00, 0x91}, []byte{0x31, 0x48})

	event, err := Parse(2400, 480, r, 0x91, ctx)
	if err != nil {
		t.Fatalf("Unexpected NoteOn event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOn event parse error - returned %v", event)
	}

	event, ok := event.(*NoteOn)
	if !ok {
		t.Fatalf("NoteOn event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid NoteOn event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestNoteOnMarshalBinary(t *testing.T) {
	evt := NoteOn{
		event: event{
			tick:    2400,
			delta:   480,
			bytes:   []byte{0x00, 0x97, 0x31, 0x48},
			tag:     types.TagNoteOn,
			Status:  0x97,
			Channel: 7,
		},
		Note: Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		},
		Velocity: 72,
	}

	expected := []byte{0x97, 0x31, 0x48}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding NoteOn (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded NoteOn\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestNoteOnUnmarshalText(t *testing.T) {
	text := "      00 97 30 48                           tick:0          delta:480        97 NoteOn                 channel:7  note:C3, velocity:72"
	expected := NoteOn{
		event: event{
			tick:    0,
			delta:   480,
			tag:     types.TagNoteOn,
			Status:  0x97,
			Channel: 7,
			bytes:   []byte{},
		},
		Note: Note{
			Value: 48,
			Name:  "C3",
			Alias: "C3",
		},
		Velocity: 72,
	}

	evt := NoteOn{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling NoteOn (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled NoteOn\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
