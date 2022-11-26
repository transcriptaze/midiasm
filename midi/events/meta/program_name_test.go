package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalProgramName(t *testing.T) {
	expected := ProgramName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   0x08,
			bytes:  []byte{0x00, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65},
		},
		Name: "Escape",
	}

	evt, err := UnmarshalProgramName(2400, 480, []byte("Escape"))
	if err != nil {
		t.Fatalf("error encoding ProgramName (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect ProgramName\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestProgramNameMarshalBinary(t *testing.T) {
	evt := ProgramName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   0x08,
			bytes:  []byte{},
		},
		Name: "Escape",
	}

	expected := []byte{0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding ProgramName (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded ProgramName\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalProgramName(t *testing.T) {
	text := "      00 FF 08 06 45 73 63 61 70 65         tick:0          delta:480        08 ProgramName            Escape"
	expected := ProgramName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   0x08,
			bytes:  []byte{},
		},
		Name: "Escape",
	}

	evt := ProgramName{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling ProgramName (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled ProgramName\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
