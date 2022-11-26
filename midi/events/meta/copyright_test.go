package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalCopyright(t *testing.T) {
	expected := Copyright{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
			bytes:  []byte{0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d},
		},
		Copyright: "Them",
	}

	evt, err := UnmarshalCopyright(2400, 480, []byte("Them"))
	if err != nil {
		t.Fatalf("error encoding Copyright (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect Copyright\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestCopyrightMarshalBinary(t *testing.T) {
	evt := Copyright{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
			bytes:  []byte{},
		},
		Copyright: "Them",
	}

	expected := []byte{0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Copyright (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Copyright\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalCopyright(t *testing.T) {
	text := "      00 FF 02 04 54 68 65 6D               tick:0          delta:480        02 Copyright              Them"
	expected := Copyright{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
			bytes:  []byte{},
		},
		Copyright: "Them",
	}

	evt := Copyright{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Copyright (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled Copyright\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestCopyrightMarshalJSON(t *testing.T) {
	tag := lib.TagCopyright

	evt := Copyright{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
			bytes:  []byte{},
		},
		Copyright: "Them",
	}

	expected := `{"tag":"Copyright","delta":480,"status":255,"type":2,"copyright":"Them"}`

	encoded, err := evt.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding %v (%v)", tag, err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded %v\n   expected:%+v\n   got:     %+v", tag, expected, string(encoded))
	}
}

func TestCopyrightUnmarshalJSON(t *testing.T) {
	text := `{"tag":"Copyright","delta":480,"status":255,"type":2,"copyright":"Them"}`
	tag := lib.TagCopyright

	expected := Copyright{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagCopyright,
			Status: 0xff,
			Type:   lib.TypeCopyright,
			bytes:  []byte{},
		},
		Copyright: "Them",
	}

	evt := Copyright{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
