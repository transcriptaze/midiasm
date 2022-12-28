package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalTempo(t *testing.T) {
	expected := Tempo{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTempo,
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20},
		},
		Tempo: 500000,
	}

	e := Tempo{}

	err := e.unmarshal(2400, 480, 0xff, []byte{0x07, 0xa1, 0x20}, []byte{0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20}...)
	if err != nil {
		t.Fatalf("error unmarshalling Tempo (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect Tempo\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestTempoMarshalBinary(t *testing.T) {
	evt := Tempo{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTempo,
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{},
		},
		Tempo: 500000,
	}

	expected := []byte{0xff, 0x51, 0x03, 0x07, 0xa1, 0x20}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Tempo (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Tempo\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTempoUnmarshalBinary(t *testing.T) {
	expected := Tempo{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagTempo,
			Status: 0xff,
			Type:   lib.TypeTempo,
			bytes:  []byte{0x83, 0x60, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20},
		},
		Tempo: 500000,
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20}

	e := Tempo{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagTempo, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagTempo, expected, e)
	}
}

func TestTempoUnmarshalText(t *testing.T) {
	text := "      00 FF 51 03 07 A1 20                  tick:0          delta:480        51 Tempo                  tempo:500000"
	expected := Tempo{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagTempo,
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{},
		},
		Tempo: 500000,
	}

	evt := Tempo{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Tempo (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled Tempo\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestTempoMarshalJSON(t *testing.T) {
	evt := Tempo{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTempo,
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{},
		},
		Tempo: 500000,
	}

	expected := `{"tag":"Tempo","delta":480,"status":255,"type":81,"tempo":500000}`

	encoded, err := evt.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding Tempo (%v)", err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded Tempo\n   expected:%+v\n   got:     %+v", expected, string(encoded))
	}
}

func TestTempoUnmarshalJSON(t *testing.T) {
	text := `{"tag":"Tempo","delta":480,"status":255,"type":81,"tempo":500000}`
	tag := lib.TagTempo

	expected := Tempo{
		event: event{
			tick:   0,
			delta:  480,
			tag:    tag,
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{},
		},
		Tempo: 500000,
	}

	evt := Tempo{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
