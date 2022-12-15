package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
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

	ctx := context.NewContext()
	e := ProgramName{}

	err := e.unmarshal(ctx, 2400, 480, 0xff, []byte("Escape"), []byte{0x00, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65}...)
	if err != nil {
		t.Fatalf("error encoding ProgramName (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect ProgramName\n   expected:%+v\n   got:     %+v", expected, e)
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

func TestProgramNameUnmarshalBinary(t *testing.T) {
	expected := ProgramName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   lib.TypeProgramName,
			bytes:  []byte{0x83, 0x60, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65},
		},
		Name: "Escape",
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65}

	e := ProgramName{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagProgramName, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagProgramName, expected, e)
	}
}

func TestProgramNameUnmarshalText(t *testing.T) {
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

func TestProgramNameMarshalJSON(t *testing.T) {
	e := ProgramName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   lib.TypeProgramName,
			bytes:  []byte{},
		},
		Name: "Escape",
	}

	expected := `{"tag":"ProgramName","delta":480,"status":255,"type":8,"name":"Escape"}`

	testMarshalJSON(t, lib.TagProgramName, e, expected)
}

func TestProgramNameNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagProgramName
	text := `{"tag":"ProgramName","delta":480,"status":255,"type":8,"name":"Escape"}`
	expected := ProgramName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagProgramName,
			Status: 0xff,
			Type:   lib.TypeProgramName,
			bytes:  []byte{},
		},
		Name: "Escape",
	}

	evt := ProgramName{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
