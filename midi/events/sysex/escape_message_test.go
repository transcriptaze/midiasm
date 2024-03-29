package sysex

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestParseEscapeMessage(t *testing.T) {
	event, err := Parse(0, false, []byte{0x83, 0x60, 0xf7, 0x02, 0xf3, 0x01}...)
	if err != nil {
		t.Fatalf("Unexpected SysEx escape message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx escape message parse error - returned %v", event)
	}

	message, ok := event.(SysExEscapeMessage)
	if !ok {
		t.Fatalf("SysEx escape message parse error - returned %T", event)
	}

	expected := lib.Hex{0xf3, 0x01}
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
			tag:    lib.TagSysExEscape,
			Status: 0xf7,
		},
		Data: lib.Hex{0xf3, 0x01},
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

func TestSysExEscapeMessageUnmarshalBinary(t *testing.T) {
	expected := SysExEscapeMessage{
		event: event{
			delta:  480,
			tag:    lib.TagSysExEscape,
			Status: 0xf7,
			bytes:  []byte{0x83, 0x60, 0xf7, 0x02, 0xf3, 0x01},
		},
		Data: lib.Hex{0xf3, 0x01},
	}

	bytes := []byte{0x83, 0x60, 0xf7, 0x02, 0xf3, 0x01}

	e := SysExEscapeMessage{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error unencoding %v (%v)", lib.TagSysExEscape, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagSysExEscape, expected, e)
	}
}

func TestSysExEscapeMessageUnmarshalText(t *testing.T) {
	text := "      00 F7 02 F3 01                        tick:1020       delta:480        F7 SysExEscape            F3 01"
	expected := SysExEscapeMessage{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSysExEscape,
			Status: 0xf7,
			bytes:  []byte{},
		},
		Data: lib.Hex{0xf3, 0x01},
	}

	evt := SysExEscapeMessage{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SysExEscapeMessage (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SysExEscapeMessage\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestSysExEscapeMessageMarshalJSON(t *testing.T) {
	e := SysExEscapeMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    lib.TagSysExEscape,
			Status: 0xf7,
		},
		Data: lib.Hex{0xf3, 0x01},
	}

	expected := `{"tag":"SysExEscape","delta":480,"status":247,"data":[243,1]}`

	testMarshalJSON(t, lib.TagSysExEscape, e, expected)
}

func TestSysExEscapeNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagSysExEscape
	text := `{"tag":"SysExEscape","delta":480,"status":247,"data":[243,1]}`
	expected := SysExEscapeMessage{
		event: event{
			tick:   0,
			delta:  480,
			bytes:  []byte{},
			tag:    lib.TagSysExEscape,
			Status: 0xf7,
		},
		Data: lib.Hex{0xf3, 0x01},
	}

	e := SysExEscapeMessage{}

	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, e)
	}
}
