package metaevent

import (
	"reflect"
	"testing"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalSequenceNumber(t *testing.T) {
	expected := SequenceNumber{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSequenceNumber,
			Status: 0xff,
			Type:   lib.TypeSequenceNumber,
			bytes:  []byte{0x00, 0xff, 0x00, 0x02, 0x00, 0x17},
		},
		SequenceNumber: 23,
	}

	evt, err := UnmarshalSequenceNumber(2400, 480, []byte{0x00, 0x17})
	if err != nil {
		t.Fatalf("error encoding SequenceNumber (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect SequenceNumber\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestSequenceNumberMarshalBinary(t *testing.T) {
	evt := SequenceNumber{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSequenceNumber,
			Status: 0xff,
			Type:   lib.TypeSequenceNumber,
			bytes:  []byte{},
		},
		SequenceNumber: 23,
	}

	expected := []byte{0xff, 0x00, 0x02, 0x00, 0x17}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SequenceNumber (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SequenceNumber\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSequenceNumberUnmarshalText(t *testing.T) {
	text := "      00 FF 00 02 00 17                     tick:0          delta:480        00 SequenceNumber         23"
	expected := SequenceNumber{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSequenceNumber,
			Status: 0xff,
			Type:   lib.TypeSequenceNumber,
			bytes:  []byte{},
		},
		SequenceNumber: 23,
	}

	evt := SequenceNumber{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SequenceNumber (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SequenceNumber\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestSequenceNumberMarshalJSON(t *testing.T) {
	tag := lib.TagSequenceNumber

	evt := SequenceNumber{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSequenceNumber,
			Status: 0xff,
			Type:   lib.TypeSequenceNumber,
			bytes:  []byte{},
		},
		SequenceNumber: 23,
	}

	expected := `{"tag":"SequenceNumber","delta":480,"status":255,"type":0,"sequence-number":23}`

	encoded, err := evt.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding %v (%v)", tag, err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded %v\n   expected:%+v\n   got:     %+v", tag, expected, string(encoded))
	}
}

func TestSequenceNumberUnmarshalJSON(t *testing.T) {
	text := `{"tag":"SequenceNumber","delta":480,"status":255,"type":0,"sequence-number":23}`
	tag := lib.TagSequenceNumber

	expected := SequenceNumber{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSequenceNumber,
			Status: 0xff,
			Type:   lib.TypeSequenceNumber,
			bytes:  []byte{},
		},
		SequenceNumber: 23,
	}

	evt := SequenceNumber{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
