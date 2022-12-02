package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParsePolyphonicPressure(t *testing.T) {
	expected := PolyphonicPressure{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xa7, 0x64},

			tag:     lib.TagPolyphonicPressure,
			Status:  0xa7,
			Channel: 7,
		},
		Pressure: 100,
	}

	ctx := context.NewContext()

	e, err := Parse(ctx, 2400, 480, 0xa7, []byte{0x64}, []byte{0x00, 0xa7, 0x64}...)
	if err != nil {
		t.Fatalf("Unexpected PolyphonicPressure event parse error: %v", err)
	} else if e == nil {
		t.Fatalf("Unexpected PolyphonicPressure event parse error - returned %v", e)
	}

	event, ok := e.(PolyphonicPressure)
	if !ok {
		t.Fatalf("PolyphonicPressure event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid PolyphonicPressure event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestPolyphonicPressureMarshalBinary(t *testing.T) {
	e := PolyphonicPressure{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xa7, 0x64},
			tag:   lib.TagPolyphonicPressure,

			Status:  0xa7,
			Channel: 7,
		},
		Pressure: 100,
	}

	expected := []byte{0xa7, 0x64}

	encoded, err := e.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding PolyphonicPressure (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded PolyphonicPressure\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestPolyphonicPressureUnmarshalText(t *testing.T) {
	text := "      00 A0 64                              tick:0          delta:480        A7 PolyphonicPressure     channel:7  pressure:100"
	expected := PolyphonicPressure{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagPolyphonicPressure,
			Status:  0xa7,
			Channel: 7,
			bytes:   []byte{},
		},
		Pressure: 100,
	}

	evt := PolyphonicPressure{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling PolyphonicPressure (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled PolyphonicPressure\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestPolyphonicPressureMarshalJSON(t *testing.T) {
	e := PolyphonicPressure{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xa7, 0x64},
			tag:   lib.TagPolyphonicPressure,

			Status:  0xa7,
			Channel: 7,
		},
		Pressure: 100,
	}

	expected := `{"tag":"PolyphonicPressure","delta":480,"status":167,"channel":7,"pressure":100}`

	testMarshalJSON(t, lib.TagPolyphonicPressure, e, expected)
}

func TestPolyphonicPressureNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagPolyphonicPressure
	text := `{"tag":"PolyphonicPressure","delta":480,"status":167,"channel":7,"pressure":100}`
	expected := PolyphonicPressure{
		event: event{
			tick:  0,
			delta: 480,
			bytes: []byte{},
			tag:   lib.TagPolyphonicPressure,

			Status:  0xa7,
			Channel: 7,
		},
		Pressure: 100,
	}

	e := PolyphonicPressure{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
