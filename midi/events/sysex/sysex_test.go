package sysex

import (
	"reflect"
	"testing"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

func TestEventMarshalBinary(t *testing.T) {
	e := event{
		tick:   0,
		delta:  480,
		bytes:  []byte{},
		tag:    lib.TagPitchBend,
		Status: lib.Status(0xf0),
	}

	expected := []byte{0x83, 0x60, 0xf0}

	if bytes, err := e.MarshalBinary(); err != nil {
		t.Fatalf("Error marshalling SYSEX base event (%v)", e)
	} else if !reflect.DeepEqual(bytes, expected) {
		t.Errorf("Incorrectly marshalled SYSEX base event\n   expected:%v\n   got:     %v", expected, bytes)
	}
}
