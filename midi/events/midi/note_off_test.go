package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParseNoteOffInMajorKey(t *testing.T) {
	expected := NoteOff{
		event{
			tick:    0,
			delta:   480,
			bytes:   []byte{0x83, 0x60, 0x81, 0x31, 0x48},
			tag:     lib.TagNoteOff,
			Status:  0x81,
			Channel: 1,
		},
		Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		}, 72,
	}

	event, err := Parse(0, 0x81, []byte{0x83, 0x60, 0x81, 0x31, 0x48}...)
	if err != nil {
		t.Fatalf("Unexpected NoteOff event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected NoteOff event parse error - returned %v", event)
	}

	event, ok := event.(NoteOff)
	if !ok {
		t.Fatalf("NoteOff event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid NoteOff event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestParseNoteOffInMinorKey(t *testing.T) {
	t.Skip()

	expected := NoteOff{
		event{
			tick:    0,
			delta:   480,
			bytes:   []byte{0x83, 0x60, 0x81, 0x00, 0x31, 0x48},
			tag:     lib.TagNoteOff,
			Status:  0x81,
			Channel: 1,
		},
		Note{
			Value: 49,
			Name:  "D♭3",
			Alias: "D♭3",
		}, 72,
	}

	event, err := Parse(0, 0x81, []byte{0x83, 0x60, 0x31, 0x48}...)
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

			tag:     lib.TagNoteOff,
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

func TestNoteOffUnmarshalBinary(t *testing.T) {
	expected := NoteOff{
		event: event{
			delta:   480,
			tag:     lib.TagNoteOff,
			Status:  0x87,
			Channel: 7,
			bytes:   []byte{0x83, 0x60, 0x87, 0x31, 0x48},
		},
		Note: Note{
			Value: 49,
			Name:  "C♯3",
			Alias: "C♯3",
		},
		Velocity: 72,
	}

	bytes := []byte{0x83, 0x60, 0x87, 0x31, 0x48}

	e := NoteOff{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagNoteOff, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagNoteOff, expected, e)
	}
}

func TestNoteOffUnmarshalText(t *testing.T) {
	text := "   83 60 87 30 40                           tick:480        delta:480        80 NoteOff                channel:7  note:C3, velocity:64"
	expected := NoteOff{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagNoteOff,
			Status:  0x87,
			Channel: 7,
			bytes:   []byte{},
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

func TestTransposeNoteOff(t *testing.T) {
	ctx := context.NewContext()

	expected := NoteOff{
		event{
			tick:    0,
			delta:   0,
			tag:     lib.TagNoteOff,
			Status:  0x81,
			Channel: 1,
		},
		Note{
			Value: 0x3a,
			Name:  "A♯3",
			Alias: "A♯3",
		}, 72,
	}

	noteOff := NoteOff{
		event{
			tick:    0,
			delta:   0,
			bytes:   []byte{0x00, 0x81, 0x39, 0x48},
			tag:     lib.TagNoteOff,
			Status:  0x81,
			Channel: 1,
		},
		Note{
			Value: 0x39,
			Name:  "A3",
			Alias: "A3",
		}, 72,
	}

	transposed := noteOff.Transpose(ctx, 1)

	if !reflect.DeepEqual(transposed, expected) {
		t.Errorf("Incorrectly transposed NoteOff event\n  expected:%#v\n  got:     %#v", expected, transposed)
	}

	if noteOff.Note.Value != 0x39 || noteOff.Note.Name != "A3" || noteOff.Note.Alias != "A3" {
		t.Errorf("Transpose mutated original NoteOff event")
	}
}

func TestNoteOffMarshalJSON(t *testing.T) {
	e := NoteOff{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagNoteOff,
			Status:  0x87,
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

	expected := `{"tag":"NoteOff","delta":480,"status":135,"channel":7,"note":{"value":48,"name":"C3","alias":"C3"},"velocity":72}`

	testMarshalJSON(t, lib.TagNoteOff, e, expected)
}

func TestNoteOffNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagNoteOff
	text := `{"tag":"NoteOff","delta":480,"status":135,"channel":7,"note":{"value":48,"name":"C3","alias":"C3"},"velocity":72}`
	expected := NoteOff{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagNoteOff,
			Status:  0x87,
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

	e := NoteOff{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
