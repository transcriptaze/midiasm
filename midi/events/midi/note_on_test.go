package midievent

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"reflect"
	"testing"
)

func TestParseNoteOnInMajorKey(t *testing.T) {
	expected := NoteOn{
		MidiEvent{
			events.Event{"NoteOn", 0, 0, 0x91, []byte{0x00, 0x91, 0x31, 0x48}},
			1,
		},
		Note{
			Value: 49,
			Name:  "C♯2",
		}, 72,
	}

	ctx := context.Context{
		Scale: context.Sharps,
	}

	e := events.Event{
		Tick:   types.Tick(0),
		Delta:  types.Delta(0),
		Status: types.Status(0x91),
		Bytes:  []byte{0x00, 0x91},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x31, 0x48}))

	event, err := Parse(e, r, &ctx)
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
		MidiEvent{
			events.Event{"NoteOn", 0, 0, 0x91, []byte{0x00, 0x91, 0x31, 0x48}},
			1,
		},
		Note{
			Value: 49,
			Name:  "D♭2",
		}, 72,
	}

	ctx := context.Context{
		Scale: context.Flats,
	}

	e := events.Event{
		Tick:   types.Tick(0),
		Delta:  types.Delta(0),
		Status: types.Status(0x91),
		Bytes:  []byte{0x00, 0x91},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x31, 0x48}))

	event, err := Parse(e, r, &ctx)
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