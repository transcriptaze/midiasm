package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestEventMarshalBinary(t *testing.T) {
	e := event{
		tick:    0,
		delta:   480,
		bytes:   []byte{},
		tag:     lib.TagPitchBend,
		Status:  lib.Status(0xe7),
		Channel: 7,
	}

	expected := []byte{0x83, 0x60, 0xe7}

	if bytes, err := e.MarshalBinary(); err != nil {
		t.Fatalf("Error marshalling MIDI base event (%v)", e)
	} else if !reflect.DeepEqual(bytes, expected) {
		t.Errorf("Incorrectly marshalled MIDI base event\n   expected:%v\n   got:     %v", expected, bytes)
	}
}
