package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestSMPTEOffsetMarshalBinary24FPS(t *testing.T) {
	evt := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        24,
		Frames:           7,
		FractionalFrames: 39,
	}

	expected := []byte{0xff, 0x54, 0x05, 0x0d, 0x2d, 0x3b, 0x07, 0x27}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SMPTE offset\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSMPTEOffsetMarshalBinary25FPS(t *testing.T) {
	evt := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
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

	expected := []byte{0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SMPTE offset\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSMPTEOffsetMarshalBinary29FPS(t *testing.T) {
	evt := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        29,
		Frames:           7,
		FractionalFrames: 39,
	}

	expected := []byte{0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SMPTE offset\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSMPTEOffsetMarshalBinary30FPS(t *testing.T) {
	evt := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        30,
		Frames:           7,
		FractionalFrames: 39,
	}

	expected := []byte{0xff, 0x54, 0x05, 0x6d, 0x2d, 0x3b, 0x07, 0x27}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SMPTE offset\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestUnmarshalSMPTEOffset24FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{0x00, 0xff, 0x54, 0x05, 0x0d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        24,
		Frames:           7,
		FractionalFrames: 39,
	}

	e := SMPTEOffset{}

	err := e.unmarshal(2400, 480, 0xff, []byte{0x0d, 0x2d, 0x3b, 0x07, 0x27}, []byte{0x00, 0xff, 0x54, 0x05, 0x0d, 0x2d, 0x3b, 0x07, 0x27}...)
	if err != nil {
		t.Fatalf("error decoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect SMPTE offset\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestUnmarshalSMPTEOffset25FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{0x00, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	}

	e := SMPTEOffset{}

	err := e.unmarshal(2400, 480, 0xff, []byte{0x2d, 0x2d, 0x3b, 0x07, 0x27}, []byte{0x00, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27}...)
	if err != nil {
		t.Fatalf("error decoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect SMPTE offset\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestUnmarshalSMPTEOffset29FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        29,
		Frames:           7,
		FractionalFrames: 39,
	}

	e := SMPTEOffset{}

	err := e.unmarshal(2400, 480, 0xff, []byte{0x4d, 0x2d, 0x3b, 0x07, 0x27}, []byte{0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27}...)
	if err != nil {
		t.Fatalf("error decoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect SMPTE offset\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestUnmarshalSMPTEOffset30FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   0x54,
			bytes:  []byte{0x00, 0xff, 0x54, 0x05, 0x6d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        30,
		Frames:           7,
		FractionalFrames: 39,
	}

	e := SMPTEOffset{}

	err := e.unmarshal(2400, 480, 0xff, []byte{0x6d, 0x2d, 0x3b, 0x07, 0x27}, []byte{0x00, 0xff, 0x54, 0x05, 0x6d, 0x2d, 0x3b, 0x07, 0x27}...)
	if err != nil {
		t.Fatalf("error decoding SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect SMPTE offset\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestSMPTEOffsetUnmarshalBinary24FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   lib.TypeSMPTEOffset,
			bytes:  []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x0d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        24,
		Frames:           7,
		FractionalFrames: 39,
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x0d, 0x2d, 0x3b, 0x07, 0x27}

	e := SMPTEOffset{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagSMPTEOffset, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagSMPTEOffset, expected, e)
	}
}

func TestSMPTEOffsetUnmarshalBinary25FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   lib.TypeSMPTEOffset,
			bytes:  []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27}

	e := SMPTEOffset{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagSMPTEOffset, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagSMPTEOffset, expected, e)
	}
}

func TestSMPTEOffsetUnmarshalBinary29FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   lib.TypeSMPTEOffset,
			bytes:  []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        29,
		Frames:           7,
		FractionalFrames: 39,
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27}

	e := SMPTEOffset{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagSMPTEOffset, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagSMPTEOffset, expected, e)
	}
}

func TestSMPTEOffsetUnmarshalBinary30FPS(t *testing.T) {
	expected := SMPTEOffset{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   lib.TypeSMPTEOffset,
			bytes:  []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x6d, 0x2d, 0x3b, 0x07, 0x27},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        30,
		Frames:           7,
		FractionalFrames: 39,
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x6d, 0x2d, 0x3b, 0x07, 0x27}

	e := SMPTEOffset{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagSMPTEOffset, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagSMPTEOffset, expected, e)
	}
}

func TestSMPTEOffsetUnmarshalText(t *testing.T) {
	text := "      00 FF 54 05 2D 2D 3B 07 27            tick:0          delta:480        54 SMPTEOffset            13 45 59 25 7 39"
	expected := SMPTEOffset{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
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
		t.Fatalf("error unmarshalling SMPTE offset (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SMPTE offset\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestSMPTEOffetMarshalJSON(t *testing.T) {
	tag := lib.TagSMPTEOffset

	evt := SMPTEOffset{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSMPTEOffset,
			Status: 0xff,
			Type:   lib.TypeSMPTEOffset,
			bytes:  []byte{},
		},
		Hour:             13,
		Minute:           45,
		Second:           59,
		FrameRate:        25,
		Frames:           7,
		FractionalFrames: 39,
	}

	expected := `{"tag":"SMPTEOffset","delta":480,"status":255,"type":84,"hour":13,"minute":45,"second":59,"frame-rate":25,"frames":7,"fractional-frames":39}`

	encoded, err := evt.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding %v (%v)", tag, err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded %v\n   expected:%+v\n   got:     %+v", tag, expected, string(encoded))
	}
}
