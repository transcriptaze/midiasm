package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParsePitchBend(t *testing.T) {
	expected := PitchBend{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xe7, 0x00, 0x08},

			tag:     lib.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
		},
		Bend: 8,
	}

	event, err := Parse(2400, 480, 0xe7, []byte{0x00, 0x08}, []byte{0x00, 0xe7, 0x00, 0x08}...)
	if err != nil {
		t.Fatalf("Unexpected PitchBend event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected PitchBend event parse error - returned %v", event)
	}

	v, ok := event.(PitchBend)
	if !ok {
		t.Fatalf("PitchBend %v type error - returned %T", event, v)
	} else if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid PitchBend event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestPitchBendMarshalBinary(t *testing.T) {
	evt := PitchBend{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xe7, 0x00, 0x08},
			tag:   lib.TagPitchBend,

			Status:  0xe7,
			Channel: 7,
		},
		Bend: 8,
	}

	expected := []byte{0xe7, 0x00, 0x08}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding PitchBend (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded PitchBend\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestPitchBendUnmarshalBinary(t *testing.T) {
	expected := PitchBend{
		event: event{
			delta:   480,
			tag:     lib.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
			bytes:   []byte{0x83, 0x60, 0xe7, 0x00, 0x08},
		},
		Bend: 8,
	}

	bytes := []byte{0x83, 0x60, 0xe7, 0x00, 0x08}

	e := PitchBend{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagPitchBend, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagPitchBend, expected, e)
	}
}

func TestPitchBendUnmarshalText(t *testing.T) {
	text := "      81 70 E7 00 08                           tick:240        delta:240        E7 PitchBend              channel:7  bend:8"
	expected := PitchBend{
		event: event{
			tick:    0,
			delta:   240,
			tag:     lib.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
			bytes:   []byte{},
		},
		Bend: 8,
	}

	evt := PitchBend{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling PitchBend (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled PitchBend\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
func TestPitchBendMarshalJSON(t *testing.T) {
	e := PitchBend{
		event: event{
			tick:    2400,
			delta:   480,
			bytes:   []byte{0x00, 0xe7, 0x00, 0x08},
			tag:     lib.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
		},
		Bend: 8,
	}

	expected := `{"tag":"PitchBend","delta":480,"status":231,"channel":7,"bend":8}`

	testMarshalJSON(t, lib.TagPitchBend, e, expected)
}

func TestPitchBendNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagPitchBend
	text := `{"tag":"PitchBend","delta":480,"status":231,"channel":7,"bend":8}`
	expected := PitchBend{
		event: event{
			tick:    0,
			delta:   480,
			bytes:   []byte{},
			tag:     lib.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
		},
		Bend: 8,
	}

	e := PitchBend{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
