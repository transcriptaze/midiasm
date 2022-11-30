package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalMIDIPort(t *testing.T) {
	expected := MIDIPort{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagMIDIPort,
			Status: 0xff,
			Type:   0x21,
			bytes:  []byte{0x00, 0xff, 0x21, 0x01, 0x70},
		},
		Port: 112,
	}

	evt, err := UnmarshalMIDIPort(2400, 480, []byte{112})
	if err != nil {
		t.Fatalf("error unmarshalling MIDIPort (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect MIDIPort\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestMIDIPortMarshalBinary(t *testing.T) {
	evt := MIDIPort{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagMIDIPort,
			Status: 0xff,
			Type:   0x21,
			bytes:  []byte{},
		},
		Port: 112,
	}

	expected := []byte{0xff, 0x21, 0x01, 0x70}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding MIDIPort (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded MIDIPort\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalMIDIPort(t *testing.T) {
	text := "      00 FF 21 01 70                        tick:0          delta:480        21 MIDIPort               112"
	expected := MIDIPort{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagMIDIPort,
			Status: 0xff,
			Type:   0x21,
			bytes:  []byte{},
		},
		Port: 112,
	}

	evt := MIDIPort{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling MIDIPort (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled MIDIPort\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestMIDIPortMarshalJSON(t *testing.T) {
	e := MIDIPort{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagMIDIPort,
			Status: 0xff,
			Type:   lib.TypeMIDIPort,
			bytes:  []byte{},
		},
		Port: 112,
	}

	expected := `{"tag":"MIDIPort","delta":480,"status":255,"type":33,"port":112}`

	testMarshalJSON(t, lib.TagMIDIPort, e, expected)
}

func TestMIDIPortNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagMIDIPort
	text := `{"tag":"MIDIPort","delta":480,"status":255,"type":33,"port":112}`
	expected := MIDIPort{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagMIDIPort,
			Status: 0xff,
			Type:   lib.TypeMIDIPort,
			bytes:  []byte{},
		},
		Port: 112,
	}

	evt := MIDIPort{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
