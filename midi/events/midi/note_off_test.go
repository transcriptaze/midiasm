package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
)

func TestParseNoteOffInMajorKey(t *testing.T) {
	expected := NoteOff{
		event{
			tick:    0,
			delta:   0,
			bytes:   []byte{0x00, 0x81, 0x31, 0x48},
			Tag:     "NoteOff",
			Status:  0x81,
			Channel: 1,
		},
		Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		}, 72,
	}

	ctx := context.NewContext()
	r := IO.TestReader([]byte{0x00, 0x81}, []byte{0x31, 0x48})

	event, err := Parse(0, 0, r, 0x81, ctx)
	if err != nil {
		t.Fatalf("Unexpected NoteOff event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOff event parse error - returned %v", event)
	}

	event, ok := event.(*NoteOff)
	if !ok {
		t.Fatalf("NoteOn event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid NoteOff event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestParseNoteOffInMinorKey(t *testing.T) {
	expected := NoteOff{
		event{
			tick:  0,
			delta: 0,
			bytes: []byte{0x00, 0x31, 0x48},

			Tag:     "NoteOff",
			Status:  0x81,
			Channel: 1,
		},
		Note{
			Value: 49,
			Name:  "D♭3",
			Alias: "D♭3",
		}, 72,
	}

	ctx := context.NewContext().UseFlats()
	r := IO.TestReader([]byte{0x00}, []byte{0x31, 0x48})

	event, err := Parse(0, 0, r, 0x81, ctx)
	if err != nil {
		t.Fatalf("Unexpected NoteOff event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOff event parse error - returned %v", event)
	}

	event, ok := event.(*NoteOff)
	if !ok {
		t.Fatalf("NoteOff event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid NoteOff event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestNoteOffMarshalBinary(t *testing.T) {
	evt := NoteOff{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0x87, 0x31, 0x48},

			Tag:     "NoteOff",
			Status:  0x87,
			Channel: 7,
		},
		Note: Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		},
		Velocity: 72,
	}

	expected := []byte{0x87, 0x31, 0x48}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding NoteOff (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded NoteOff\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestNoteOffUnmarshalText(t *testing.T) {
	text := "   83 60 87 30 40                           tick:480        delta:480        80 NoteOff                channel:7  note:C3, velocity:64"
	expected := NoteOff{
		event: event{
			tick:    0,
			delta:   480,
			Tag:     "NoteOff",
			Status:  0x87,
			Channel: 7,
			bytes:   []byte{0x00, 0x87, 0x30, 0x40},
		},
		Note: Note{
			Value: 48,
			Name:  "C3",
			Alias: "C3",
		},
		Velocity: 64,
	}

	evt := NoteOff{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling NoteOff (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled NoteOff\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
