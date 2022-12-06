package sysex

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type TestSysExEvent interface {
	SysExMessage | SysExContinuationMessage | SysExEscapeMessage

	MarshalJSON() ([]byte, error)
}

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

func testMarshalJSON[E TestSysExEvent](t *testing.T, tag lib.Tag, e E, expected string) {
	encoded, err := e.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding %v (%v)", tag, err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded %v\n   expected:%+v\n   got:     %+v", tag, expected, string(encoded))
	}
}
