package midievent

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
)

func TestParseNoteOnInMajorKey(t *testing.T) {
	expected := NoteOn{
		"NoteOn",
		0x91,
		1,
		Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		}, 72,
	}

	ctx := context.NewContext()
	r := bufio.NewReader(bytes.NewReader([]byte{0x31, 0x48}))

	event, err := Parse(0, 0, r, 0x91, ctx)
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
		"NoteOn",
		0x91,
		1,
		Note{
			Value: 49,
			Name:  "D♭3",
			Alias: "D♭3",
		}, 72,
	}

	ctx := context.NewContext().UseFlats()
	r := bufio.NewReader(bytes.NewReader([]byte{0x31, 0x48}))

	event, err := Parse(0, 0, r, 0x91, ctx)
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
