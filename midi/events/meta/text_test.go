package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalText(t *testing.T) {
	expected := Text{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagText,
			Status: 0xff,
			Type:   0x01,
			bytes:  []byte{0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74},
		},
		Text: "This and That",
	}

	evt, err := UnmarshalText(2400, 480, []byte("This and That"))
	if err != nil {
		t.Fatalf("error encoding Text (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect Text\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestTextMarshalBinary(t *testing.T) {
	evt := Text{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagText,
			Status: 0xff,
			Type:   0x01,
			bytes:  []byte{},
		},
		Text: "This and That",
	}

	expected := []byte{0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Text (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Text\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalText(t *testing.T) {
	text := "      00 FF 01 0D 54 68 69 73 20 61 6E 64…  tick:0          delta:480        01 Text                   This and That"
	expected := Text{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagText,
			Status: 0xff,
			Type:   0x01,
			bytes:  []byte{},
		},
		Text: "This and That",
	}

	evt := Text{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Text (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled Text\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
