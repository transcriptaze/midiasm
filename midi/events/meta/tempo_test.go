package metaevent

import (
	"reflect"
	"testing"
)

func TestTempoMarshalBinary(t *testing.T) {
	tempo := Tempo{
		event: event{
			tick:   2400,
			delta:  480,
			Tag:    "Tempo",
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{},
		},
		Tempo: 500000,
	}

	expected := []byte{0xff, 0x51, 0x03, 0x07, 0xa1, 0x20}

	encoded, err := tempo.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Tempo (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Tempo\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTempoUnmarshalText(t *testing.T) {
	text := "      00 FF 51 03 07 A1 20                  tick:0          delta:0          51 Tempo                  tempo:500000"
	expected := Tempo{
		event: event{
			tick:   0,
			delta:  0,
			Tag:    "Tempo",
			Status: 0xff,
			Type:   0x51,
			bytes:  []byte{},
		},
		Tempo: 500000,
	}

	tempo := Tempo{}

	if err := tempo.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Tempo (%v)", err)
	}

	if !reflect.DeepEqual(tempo, expected) {
		t.Errorf("incorrectly unmarshalled Tempo\n   expected:%+v\n   got:     %+v", expected, tempo)
	}

}
