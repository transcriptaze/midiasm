package sysex

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParseEscapeMessage(t *testing.T) {
	ctx := context.NewContext()
	r := IO.TestReader([]byte{}, []byte{0x02, 0xf3, 0x01})

	event, err := Parse(0, 0, r, 0xf7, ctx)
	if err != nil {
		t.Fatalf("Unexpected SysEx escape message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx escape message parse error - returned %v", event)
	}

	message, ok := event.(*SysExEscapeMessage)
	if !ok {
		t.Fatalf("SysEx escape message parse error - returned %T", event)
	}

	expected := types.Hex{0xf3, 0x01}
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx escape message data - expected:%v, got: %v", expected, message.Data)
	}
}

func TestSysExEscapeMessageMarshalBinary(t *testing.T) {
	evt := SysExEscapeMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    types.TagSysExEscape,
			Status: 0xf7,
		},
		Data: types.Hex{0xf3, 0x01},
	}

	expected := []byte{0xf7, 0x02, 0xf3, 0x01}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SysExEscapeMessage (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SysExEscapeMessage\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSysExEscapeMessageUnmarshalText(t *testing.T) {
	text := "      00 F7 02 F3 01                        tick:1020       delta:480        F7 SysExEscape            F3 01"
	expected := SysExEscapeMessage{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagSysExEscape,
			Status: 0xf7,
			bytes:  []byte{},
		},
		Data: types.Hex{0xf3, 0x01},
	}

	evt := SysExEscapeMessage{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SysExEscapeMessage (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SysExEscapeMessage\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
