package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestProgramChange(t *testing.T) {
	expected := ProgramChange{
		event: event{
			tick:    12345,
			delta:   5432,
			bytes:   []byte{0x00, 0xc7, 0x0d},
			tag:     lib.TagProgramChange,
			Status:  0xc7,
			Channel: 7,
		},
		Bank:    0,
		Program: 13,
	}

	ctx := context.NewContext()
	ctx.ProgramBank[7] = 673

	event, err := Parse(ctx, 12345, 5432, 0xc7, []byte{0x0d}, []byte{0x00, 0xc7, 0x0d}...)
	if err != nil {
		t.Fatalf("Unexpected ProgramChange event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected ProgramChange event parse error - returned %v", event)
	}

	event, ok := event.(ProgramChange)
	if !ok {
		t.Fatalf("ProgramChange event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid ProgramChange event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestProgramChangeMarshalBinary(t *testing.T) {
	evt := ProgramChange{
		event: event{
			tick:    2400,
			delta:   480,
			bytes:   []byte{0x00, 0xc7, 25},
			tag:     lib.TagProgramChange,
			Status:  0xc7,
			Channel: 7,
		},
		Bank:    0,
		Program: 25,
	}

	expected := []byte{0xc7, 0x19}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding ProgramChange (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded ProgramChange\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestProgramChangeUnmarshalBinary(t *testing.T) {
	expected := ProgramChange{
		event: event{
			delta:   480,
			tag:     lib.TagProgramChange,
			Status:  0xc7,
			Channel: 7,
			bytes:   []byte{0x83, 0x60, 0xc7, 0x19},
		},
		Bank:    0,
		Program: 25,
	}

	bytes := []byte{0x83, 0x60, 0xc7, 0x19}

	e := ProgramChange{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagProgramChange, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagProgramChange, expected, e)
	}
}

func TestProgramChangeUnmarshalText(t *testing.T) {
	text := "      00 C7 19                              tick:0          delta:480        C7 ProgramChange          channel:7  bank:3, program:25"
	expected := ProgramChange{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagProgramChange,
			Status:  0xc7,
			Channel: 7,
			bytes:   []byte{},
		},
		Bank:    3,
		Program: 25,
	}

	evt := ProgramChange{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling ProgramChange (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled ProgramChange\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestProgramChangeMarshalJSON(t *testing.T) {
	e := ProgramChange{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xb7, 0x54, 0x1d},
			tag:   lib.TagProgramChange,

			Status:  0xc7,
			Channel: 7,
		},
		Bank:    673,
		Program: 13,
	}

	expected := `{"tag":"ProgramChange","delta":480,"status":199,"channel":7,"bank":673,"program":13}`

	testMarshalJSON(t, lib.TagProgramChange, e, expected)
}

func TestProgramChangeNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagProgramChange
	text := `{"tag":"ProgramChange","delta":480,"status":199,"channel":7,"bank":673,"program":13}`
	expected := ProgramChange{
		event: event{
			tick:  0,
			delta: 480,
			bytes: []byte{},
			tag:   lib.TagProgramChange,

			Status:  0xc7,
			Channel: 7,
		},
		Bank:    673,
		Program: 13,
	}

	e := ProgramChange{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
