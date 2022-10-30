package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalCuePoint(t *testing.T) {
	expected := CuePoint{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagCuePoint,
			Status: 0xff,
			Type:   0x07,
			bytes:  []byte{0x00, 0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c},
		},
		CuePoint: "More cowbell",
	}

	evt, err := UnmarshalCuePoint(2400, 480, []byte("More cowbell"))
	if err != nil {
		t.Fatalf("error encoding CuePoint (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect CuePoint\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestCuePointMarshalBinary(t *testing.T) {
	evt := CuePoint{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagCuePoint,
			Status: 0xff,
			Type:   0x07,
			bytes:  []byte{},
		},
		CuePoint: "More cowbell",
	}

	expected := []byte{0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding CuePoint (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded CuePoint\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalCuePoint(t *testing.T) {
	text := "      00 FF 07 0C 4D 6F 72 65 20 63 6F 77…  tick:0          delta:480        07 CuePoint               More cowbell"
	expected := CuePoint{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagCuePoint,
			Status: 0xff,
			Type:   0x07,
			bytes:  []byte{},
		},
		CuePoint: "More cowbell",
	}

	evt := CuePoint{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling CuePoint (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled CuePoint\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
