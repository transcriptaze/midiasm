package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalMIDIChannelPrefix(t *testing.T) {
	expected := MIDIChannelPrefix{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagMIDIChannelPrefix,
			Status: 0xff,
			Type:   0x20,
			bytes:  []byte{0x00, 0xff, 0x20, 0x01, 0x0d},
		},
		Channel: 13,
	}

	evt, err := UnmarshalMIDIChannelPrefix(2400, 480, []byte{13})
	if err != nil {
		t.Fatalf("error unmarshalling MIDIChannelPrefix (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect MIDIChannelPrefix\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestMIDIChannelPrefixMarshalBinary(t *testing.T) {
	evt := MIDIChannelPrefix{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    types.TagMIDIChannelPrefix,
			Status: 0xff,
			Type:   0x20,
			bytes:  []byte{},
		},
		Channel: 13,
	}

	expected := []byte{0xff, 0x20, 0x01, 0x0d}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding MIDIChannelPrefix (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded MIDIChannelPrefix\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalMIDIChannelPrefix(t *testing.T) {
	text := "      00 FF 20 01 0D                        tick:0          delta:480        20 MIDIChannelPrefix      13"
	expected := MIDIChannelPrefix{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagMIDIChannelPrefix,
			Status: 0xff,
			Type:   0x20,
			bytes:  []byte{},
		},
		Channel: 13,
	}

	evt := MIDIChannelPrefix{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling MIDIChannelPrefix (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled MIDIChannelPrefix\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
