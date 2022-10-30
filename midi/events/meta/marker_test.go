package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalMarker(t *testing.T) {
	expected := Marker{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagMarker,
			Status: 0xff,
			Type:   0x06,
			bytes:  []byte{0x00, 0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73},
		},
		Marker: "Here Be Dragons",
	}

	evt, err := UnmarshalMarker(2400, 480, []byte("Here Be Dragons"))
	if err != nil {
		t.Fatalf("error encoding Marker (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect Marker\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestMarkerMarshalBinary(t *testing.T) {
	evt := Marker{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagMarker,
			Status: 0xff,
			Type:   0x06,
			bytes:  []byte{},
		},
		Marker: "Here Be Dragons",
	}

	expected := []byte{0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Marker (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Marker\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalMarker(t *testing.T) {
	text := "      00 FF 06 0F 48 65 72 65 20 42 65 20…  tick:0          delta:480        06 Marker                 Here Be Dragons"
	expected := Marker{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagMarker,
			Status: 0xff,
			Type:   0x06,
			bytes:  []byte{},
		},
		Marker: "Here Be Dragons",
	}

	evt := Marker{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Marker (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled Marker\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
