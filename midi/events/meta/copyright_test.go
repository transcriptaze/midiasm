package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalCopyright(t *testing.T) {
	expected := Copyright{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagCopyright,
			Status: 0xff,
			Type:   0x02,
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
			tag:    types.TagCopyright,
			Status: 0xff,
			Type:   0x02,
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
			tag:    types.TagCopyright,
			Status: 0xff,
			Type:   0x02,
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
