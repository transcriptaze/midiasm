package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalCMajorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  nil,
		},
		Accidentals: 0,
		KeyType:     lib.Major,
		Key:         "C major",
	}

	event, err := UnmarshalKeySignature(2400, 480, []byte{0x00, 0x00}...)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestUnmarshalCMinorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  nil,
		},
		Accidentals: -3,
		KeyType:     lib.Minor,
		Key:         "C minor",
	}

	event, err := UnmarshalKeySignature(2400, 480, []byte{0xfd, 0x01}...)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestKeySignatureMarshalBinary(t *testing.T) {
	evt := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{},
		},
		Accidentals: -3,
		KeyType:     lib.Minor,
		Key:         "C minor",
	}

	expected := []byte{0xff, 0x59, 0x02, 0xfd, 0x01}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding KeySignature (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded KeySignature\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestKeySignatureUnmarshalText(t *testing.T) {
	text := "      00 FF 59 02 00 01                     tick:0          delta:480        59 KeySignature           B minor"
	expected := KeySignature{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{},
		},
		Accidentals: 2,
		KeyType:     1,
		Key:         "B minor",
	}

	evt := KeySignature{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling KeySignature (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled KeySignature\n   expected:%#v\n   got:     %#v", expected, evt)
	}

}
