package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalSMPTEOffset(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	}

	evt, err := UnmarshalSMPTEOffset(2400, 480, []byte{0x4d, 0x2d, 0x3b, 0x07, 0x27})
	if err != nil {
		t.Fatalf("error encoding SMPTEOffset (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect SMPTEOffset\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestSMPTEMarshalBinary(t *testing.T) {
	evt := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	}

	expected := []byte{0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SMPTEOffset (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SMPTEOffset\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSMPTEOffsetUnmarshalText(t *testing.T) {
	text := "      00 FF 54 05 4D 2D 3B 07 27            tick:0          delta:480        54 SMPTEOffset            13 45 59 25 7 39"
	expected := SMPTEOffset{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	}

	evt := SMPTEOffset{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SMPTEOffset (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SMPTEOffset\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
