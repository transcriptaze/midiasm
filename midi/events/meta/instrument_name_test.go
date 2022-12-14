package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalInstrumentName(t *testing.T) {
	expected := InstrumentName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
			bytes:  []byte{0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
		},
		Name: "Didgeridoo",
	}

	ctx := context.NewContext()
	e := InstrumentName{}

	err := e.unmarshal(ctx, 2400, 480, 0xff, []byte("Didgeridoo"), []byte{0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f}...)
	if err != nil {
		t.Fatalf("error encoding InstrumentName (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect InstrumentName\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestInstrumentNameMarshalBinary(t *testing.T) {
	evt := InstrumentName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
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

func TestInstrumentnameUnmarshalBinary(t *testing.T) {
	expected := InstrumentName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
			bytes:  []byte{0x83, 0x60, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
		},
		Name: "Didgeridoo",
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f}

	e := InstrumentName{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagInstrumentName, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagInstrumentName, expected, e)
	}
}

func TestTextUnmarshalInstrumentName(t *testing.T) {
	text := "      00 FF 04 0A 44 69 64 67 65 72 69 64…  tick:0          delta:480        04 InstrumentName         Didgeridoo"
	expected := InstrumentName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
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

func TestInstrumentNameMarshalJSON(t *testing.T) {
	e := InstrumentName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
			bytes:  []byte{},
		},
		Name: "Didgeridoo",
	}

	expected := `{"tag":"InstrumentName","delta":480,"status":255,"type":4,"name":"Didgeridoo"}`

	testMarshalJSON(t, lib.TagInstrumentName, e, expected)
}

func TestInstrumentNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagInstrumentName
	text := `{"tag":"InstrumentName","delta":480,"status":255,"type":4,"name":"Didgeridoo"}`
	expected := InstrumentName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagInstrumentName,
			Status: 0xff,
			Type:   lib.TypeInstrumentName,
			bytes:  []byte{},
		},
		Name: "Didgeridoo",
	}

	evt := InstrumentName{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
