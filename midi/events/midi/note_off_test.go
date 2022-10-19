package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
)

func TestParseNoteOffInMajorKey(t *testing.T) {
	expected := NoteOff{
		"NoteOff",
		0x81,
		1,
		Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		}, 72,
	}

	ctx := context.NewContext()
	r := IO.BytesReader([]byte{0x31, 0x48})

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
		"NoteOff",
		0x81,
		1,
		Note{
			Value: 49,
			Name:  "D♭3",
			Alias: "D♭3",
		}, 72,
	}

	ctx := context.NewContext().UseFlats()
	r := IO.BytesReader([]byte{0x31, 0x48})

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
