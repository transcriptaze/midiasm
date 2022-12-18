package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParseChannelPressure(t *testing.T) {
	expected := ChannelPressure{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xd7, 0x64},

			tag:     lib.TagChannelPressure,
			Status:  0xd7,
			Channel: 7,
		},
		Pressure: 100,
	}

	ctx := context.NewContext()

	event, err := Parse(ctx, 2400, 480, 0xd7, []byte{0x64}, []byte{0x00, 0xd7, 0x64}...)
	if err != nil {
		t.Fatalf("Unexpected ChannelPressure event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected ChannelPressure event parse error - returned %v", event)
	}

	event, ok := event.(ChannelPressure)
	if !ok {
		t.Fatalf("ChannelPressure event parse error - returned %T", event)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("Invalid ChannelPressure event\n  expected:%#v\n  got:     %#v", expected, event)
	}
}

func TestChannelPressureMarshalBinary(t *testing.T) {
	evt := ChannelPressure{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xd7, 0x64},
			tag:   lib.TagChannelPressure,

			Status:  0xd7,
			Channel: 7,
		},
		Pressure: 100,
	}

	expected := []byte{0xd7, 0x64}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding ChannelPressure (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded ChannelPressure\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestChannelPressureUnmarshalBinary(t *testing.T) {
	expected := ChannelPressure{
		event: event{
			delta:   480,
			tag:     lib.TagChannelPressure,
			Status:  0xd7,
			Channel: 7,
			bytes:   []byte{0x83, 0x60, 0xd7, 0x64},
		},
		Pressure: 100,
	}

	bytes := []byte{0x83, 0x60, 0xd7, 0x64}

	e := ChannelPressure{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagChannelPressure, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagChannelPressure, expected, e)
	}
}

func TestChannelPressureUnmarshalText(t *testing.T) {
	text := "      00 D0 07                              tick:0          delta:480        D0 ChannelPressure        channel:7  pressure:100"
	expected := ChannelPressure{
		event: event{
			tick:    0,
			delta:   480,
			tag:     lib.TagChannelPressure,
			Status:  0xd7,
			Channel: 7,
			bytes:   []byte{},
		},
		Pressure: 100,
	}

	evt := ChannelPressure{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling ChannelPressure (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled ChannelPressure\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestChannelPressureMarshalJSON(t *testing.T) {
	e := ChannelPressure{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xa7, 0x64},
			tag:   lib.TagChannelPressure,

			Status:  0xd7,
			Channel: 7,
		},
		Pressure: 100,
	}

	expected := `{"tag":"ChannelPressure","delta":480,"status":215,"channel":7,"pressure":100}`

	testMarshalJSON(t, lib.TagChannelPressure, e, expected)
}

func TestChannelPressureNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagChannelPressure
	text := `{"tag":"ChannelPressure","delta":480,"status":215,"channel":7,"pressure":100}`
	expected := ChannelPressure{
		event: event{
			tick:  0,
			delta: 480,
			bytes: []byte{},
			tag:   lib.TagChannelPressure,

			Status:  0xd7,
			Channel: 7,
		},
		Pressure: 100,
	}

	e := ChannelPressure{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
