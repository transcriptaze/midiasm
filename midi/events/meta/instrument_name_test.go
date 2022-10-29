package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalInstrumentName(t *testing.T) {
	expected := InstrumentName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagInstrumentName,
			Status: 0xff,
			Type:   0x04,
			bytes:  []byte{0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
		},
		Name: "Didgeridoo",
	}

	evt, err := UnmarshalInstrumentName(2400, 480, []byte("Didgeridoo"))
	if err != nil {
		t.Fatalf("error encoding InstrumentName (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect InstrumentName\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestInstrumentNameMarshalBinary(t *testing.T) {
	evt := InstrumentName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagInstrumentName,
			Status: 0xff,
			Type:   0x04,
			bytes:  []byte{},
		},
		Name: "Didgeridoo",
	}

	expected := []byte{0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding InstrumentName (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded InstrumentName\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalInstrumentName(t *testing.T) {
	text := "      00 FF 04 0A 44 69 64 67 65 72 69 64…  tick:0          delta:480        04 InstrumentName         Didgeridoo"
	expected := InstrumentName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagInstrumentName,
			Status: 0xff,
			Type:   0x04,
			bytes:  []byte{},
		},
		Name: "Didgeridoo",
	}

	evt := InstrumentName{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling InstrumentName (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled InstrumentName\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
