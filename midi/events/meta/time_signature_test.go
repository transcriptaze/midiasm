package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalTimeSignature(t *testing.T) {
	expected := TimeSignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTimeSignature,
			Status: 0xff,
			Type:   lib.TypeTimeSignature,
			bytes:  []byte{0x00, 0xff, 0x58, 0x04, 0x03, 0x02, 0x18, 0x08},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	e := TimeSignature{}

	err := e.unmarshal(2400, 480, 0xff, []byte{0x03, 0x02, 0x18, 0x08}, []byte{0x00, 0xff, 0x58, 0x04, 0x03, 0x02, 0x18, 0x08}...)
	if err != nil {
		t.Fatalf("error unmarshalling TimeSignature (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect TimeSignature\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestTimeSignatureMarshalBinary(t *testing.T) {
	evt := TimeSignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTimeSignature,
			Status: 0xff,
			Type:   lib.TypeTimeSignature,
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

func TestTimeSignatureUnmarshalBinary(t *testing.T) {
	expected := TimeSignature{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagTimeSignature,
			Status: 0xff,
			Type:   lib.TypeTimeSignature,
			bytes:  []byte{0x83, 0x60, 0xff, 0x58, 0x04, 0x03, 0x02, 0x18, 0x08},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x58, 0x04, 0x03, 0x02, 0x18, 0x08}

	e := TimeSignature{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagTimeSignature, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagTimeSignature, expected, e)
	}
}

func TestTimeSignatureUnmarshalText(t *testing.T) {
	text := "      00 FF 58 04 04 02 18 08               tick:0          delta:480        58 TimeSignature          3/4, 24 ticks per click, 8/32 per quarter"
	tag := lib.TagTimeSignature

	expected := TimeSignature{
		event: event{
			tick:   0,
			delta:  480,
			tag:    tag,
			Status: 0xff,
			Type:   lib.TypeTimeSignature,
			bytes:  []byte{},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	evt := TimeSignature{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%#v\n   got:     %#v", tag, expected, evt)
	}

}

func TestTimeSignatureMarshalJSON(t *testing.T) {
	tag := lib.TagTimeSignature

	evt := TimeSignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    tag,
			Status: 0xff,
			Type:   0x58,
			bytes:  []byte{},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	expected := `{"tag":"TimeSignature","delta":480,"status":255,"type":88,"numerator":3,"denominator":4,"ticks-per-click":24,"thirty-seconds-per-quarter":8}`

	encoded, err := evt.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding %v (%v)", tag, err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded %v\n   expected:%+v\n   got:     %+v", tag, expected, string(encoded))
	}
}

func TestTimeSignatureUnmarshalJSON(t *testing.T) {
	text := `{"tag":"TimeSignature","delta":480,"status":255,"type":88,"numerator":3,"denominator":4,"ticks-per-click":24,"thirty-seconds-per-quarter":8}`
	tag := lib.TagTimeSignature

	expected := TimeSignature{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagTimeSignature,
			Status: 0xff,
			Type:   lib.TypeTimeSignature,
			bytes:  []byte{},
		},
		Numerator:               3,
		Denominator:             4,
		TicksPerClick:           24,
		ThirtySecondsPerQuarter: 8,
	}

	evt := TimeSignature{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
