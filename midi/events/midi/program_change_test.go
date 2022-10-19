package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
)

func TestProgramChange(t *testing.T) {
	expected := ProgramChange{
		event: event{
			tick:  12345,
			delta: 5432,
			bytes: []byte{0x00, 0xc7, 0x0d},

			Tag:     "ProgramChange",
			Status:  0xc7,
			Channel: 7,
		},
		Bank:    673,
		Program: 13,
	}

	r := IO.BytesReader([]byte{0x0d})

	ctx := context.NewContext()
	ctx.ProgramBank[7] = 673

	event, err := Parse(12345, 5432, r, 0xc7, ctx)
	if err != nil {
		t.Fatalf("Unexpected ProgramChange event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected ProgramChange event parse error - returned %v", event)
	}

	event, ok := event.(*ProgramChange)
	if !ok {
		t.Fatalf("ProgramChange event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid ProgramChange event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestProgramChangeMarshalBinary(t *testing.T) {
	evt := ProgramChange{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xc7, 25},

			Tag:     "ProgramChange",
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

func TestProgramChangeUnmarshalText(t *testing.T) {
	text := "      00 C7 19                              tick:0          delta:0          C7 ProgramChange          channel:7  bank:3, program:25"
	expected := ProgramChange{
		event: event{
			tick:    0,
			delta:   0,
			Tag:     "ProgramChange",
			Status:  0xc7,
			Channel: 7,
			bytes:   []byte{0x00, 0xc7, 0x19},
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
