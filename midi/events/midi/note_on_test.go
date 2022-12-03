package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParseNoteOnInMajorKey(t *testing.T) {
	expected := NoteOn{
		event: event{
			tick:    2400,
			delta:   480,
			bytes:   []byte{0x00, 0x91, 0x31, 0x48},
			tag:     lib.TagNoteOn,
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

	event, err := Parse(ctx, 2400, 480, 0x91, []byte{0x31, 0x48}, []byte{0x00, 0x91, 0x31, 0x48}...)
	if err != nil {
		t.Fatalf("Unexpected NoteOn event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOn event parse error - returned %v", event)
	}

	event, ok := event.(NoteOn)
	if !ok {
		t.Fatalf("NoteOn event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid NoteOn event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestParseNoteOnInMinorKey(t *testing.T) {
	expected := NoteOn{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0x91, 0x31, 0x48},

			tag:     lib.TagNoteOn,
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

	event, err := Parse(ctx, 2400, 480, 0x91, []byte{0x31, 0x48}, []byte{0x00, 0x91, 0x31, 0x48}...)
	if err != nil {
		t.Fatalf("Unexpected NoteOn event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOn event parse error - returned %v", event)
	}

	event, ok := event.(NoteOn)
	if !ok {
		t.Fatalf("NoteOn event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid NoteOn event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestNoteOnMarshalBinary(t *testing.T) {
	evt := NoteOn{
		event: event{
			tick:    2400,
			delta:   480,
			bytes:   []byte{0x00, 0x97, 0x31, 0x48},
			tag:     lib.TagNoteOn,
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
			tag:     lib.TagNoteOn,
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

func TestTransposeNoteOn(t *testing.T) {
	ctx := context.NewContext()

	expected := NoteOn{
		event{
			tick:    0,
			delta:   0,
			tag:     lib.TagNoteOn,
			Status:  0x91,
			Channel: 1,
		},
		Note{
			Value: 0x3a,
			Name:  "A♯3",
			Alias: "A♯3",
		}, 72,
	}

	noteOn := NoteOn{
		event{
			tick:    0,
			delta:   0,
			bytes:   []byte{0x00, 0x81, 0x39, 0x48},
			tag:     lib.TagNoteOn,
			Status:  0x91,
			Channel: 1,
		},
		Note{
			Value: 0x39,
			Name:  "A3",
			Alias: "A3",
		}, 72,
	}

	transposed := noteOn.Transpose(ctx, 1)

	if !reflect.DeepEqual(transposed, expected) {
		t.Errorf("Incorrectly transposed NoteOn event\n  expected:%#v\n  got:     %#v", expected, transposed)
	}

	if noteOn.Note.Value != 0x39 || noteOn.Note.Name != "A3" || noteOn.Note.Alias != "A3" {
		t.Errorf("Transpose mutated original NoteOn event")
	}
}

func TestNoteOnMarshalJSON(t *testing.T) {
	e := NoteOn{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagNoteOn,
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

	expected := `{"tag":"NoteOn","delta":480,"status":151,"channel":7,"note":{"value":48,"name":"C3","alias":"C3"},"velocity":72}`

	testMarshalJSON(t, lib.TagNoteOn, e, expected)
}

func TestNoteOnNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagNoteOn
	text := `{"tag":"NoteOn","delta":480,"status":151,"channel":7,"note":{"value":48,"name":"C3","alias":"C3"},"velocity":72}`
	expected := NoteOn{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagNoteOn,
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

	e := NoteOn{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
