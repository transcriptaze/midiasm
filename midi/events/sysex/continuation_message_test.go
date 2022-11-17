package sysex

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParseContinuationMessage(t *testing.T) {
	ctx := context.NewContext()
	ctx.Casio = true
	bytes := []byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}

	event, err := Parse(ctx, 0, 0, 0xf7, bytes[1:], bytes...)
	if err != nil {
		t.Fatalf("Unexpected SysEx continuation message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx continuation message parse error - returned %v", event)
	}

	message, ok := event.(*SysExContinuationMessage)
	if !ok {
		t.Fatalf("SysEx continuation message parse error - returned %T", event)
	}

	expected := types.Hex{0x7e, 0x00, 0x09, 0x01}
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx continuation message data - expected:%v, got: %v", expected, message.Data)
	}

	if ctx.Casio {
		t.Errorf("context Casio flag not reset")
	}
}

func TestSysExContinuationMessageMarshalBinary(t *testing.T) {
	evt := SysExContinuationMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    types.TagSysExContinuation,
			Status: 0xf7,
		},
		Data: types.Hex{0x7e, 0x00, 0x09, 0x01},
	}

	expected := []byte{0xf7, 0x04, 0x7e, 0x00, 0x09, 0x01}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SysExContinuationMessage (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SysExContinuationMessage\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSysExContinuationEndMessageMarshalBinary(t *testing.T) {
	evt := SysExContinuationMessage{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{},
			tag:    types.TagSysExContinuation,
			Status: 0xf7,
		},
		Data: types.Hex{0x7e, 0x00, 0x09, 0x01},
		End:  true,
	}

	expected := []byte{0xf7, 0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SysExContinuationEndMessage (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SysExContinuationEndMessage\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestSysExContinuationMessageUnmarshalText(t *testing.T) {
	text := "   81 48 F7 06 43 12 00 43 12 00            tick:920        delta:200        F7 SysExContinuation      43 12 00 43 12 00"
	expected := SysExContinuationMessage{
		event: event{
			tick:   0,
			delta:  200,
			tag:    types.TagSysExContinuation,
			Status: 0xf7,
			bytes:  []byte{},
		},
		Data: types.Hex{0x43, 0x12, 0x00, 0x43, 0x12, 0x00},
	}

	evt := SysExContinuationMessage{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SysExContinuationMessage (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SysExContinuationMessage\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}
