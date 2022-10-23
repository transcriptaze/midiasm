package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParseTimeSignature(t *testing.T) {
	event, err := NewTimeSignature(2400, 480, []byte{3, 3, 24, 8})
	if err != nil {
		t.Fatalf("TimeSignature parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("TimeSignature parse error - returned %v", event)
	}

	if event.Numerator != 3 {
		t.Errorf("Invalid TimeSignature numerator - expected:%v, got: %v", 3, event.Numerator)
	}

	if event.Denominator != 8 {
		t.Errorf("Invalid TimeSignature denominator - expected:%v, got: %v", 8, event.Denominator)
	}
}

func TestTimeSignatureMarshalBinary(t *testing.T) {
	evt := TimeSignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagTimeSignature,
			Status: 0xff,
			Type:   0x58,
			bytes:  []byte{},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	expected := []byte{0xff, 0x58, 0x04, 0x03, 0x02, 0x18, 0x08}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding TimeSignature (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded TimeSignature\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTimeSignatureUnmarshalText(t *testing.T) {
	text := "      00 FF 58 04 04 02 18 08               tick:0          delta:480        58 TimeSignature          3/4, 24 ticks per click, 8/32 per quarter"
	expected := TimeSignature{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagTimeSignature,
			Status: 0xff,
			Type:   0x58,
			bytes:  []byte{},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	evt := TimeSignature{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling TimeSignature (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled TimeSignature\n   expected:%#v\n   got:     %#v", expected, evt)
	}

}
