package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
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
			bytes:  []byte{0xff, 0x59, 0x02, 0x00, 0x00},
		},
		Accidentals: 0,
		KeyType:     lib.Major,
		Key:         "C major",
	}

	ctx := context.NewContext()
	e := KeySignature{}

	err := e.unmarshal(ctx, 2400, 480, 0xff, []byte{0x00, 0x00}, []byte{0xff, 0x59, 0x02, 0x00, 0x00}...)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", expected, e)
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
			bytes:  []byte{0xff, 0x59, 0x02, 0xfd, 0x01},
		},
		Accidentals: -3,
		KeyType:     lib.Minor,
		Key:         "C minor",
	}

	ctx := context.NewContext()
	e := KeySignature{}

	err := e.unmarshal(ctx, 2400, 480, 0xff, []byte{0xfd, 0x01}, []byte{0xff, 0x59, 0x02, 0xfd, 0x01}...)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", expected, e)
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

func TestKeySignatureMarshalJSON(t *testing.T) {
	e := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   lib.TypeKeySignature,
			bytes:  []byte{},
		},
		Accidentals: -3,
		KeyType:     lib.Minor,
		Key:         "C minor",
	}

	expected := `{"tag":"KeySignature","delta":480,"status":255,"type":89,"accidentals":-3,"key-type":1,"key":"C minor"}`

	testMarshalJSON(t, lib.TagKeySignature, e, expected)
}

func TestKeySignatureNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagKeySignature
	text := `{"tag":"KeySignature","delta":480,"status":255,"type":89,"accidentals":-3,"key-type":1,"key":"C minor"}`
	expected := KeySignature{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   lib.TypeKeySignature,
			bytes:  []byte{},
		},
		Accidentals: -3,
		KeyType:     lib.Minor,
		Key:         "C minor",
	}

	evt := KeySignature{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
