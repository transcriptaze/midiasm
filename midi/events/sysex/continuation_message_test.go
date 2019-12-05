package sysex

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"reflect"
	"testing"
)

func TestParseContinuationMessage(t *testing.T) {
	e := events.Event{
		Tick:   0,
		Delta:  0,
		Status: 0xf7,
		Bytes:  []byte{0x00, 0xf7},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}))

	event, err := Parse(e, r)
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

	expected := []byte{0x7e, 0x00, 0x09, 0x01, 0xf7}
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx continuation message data - expected:%v, got: %v", expected, message.Data)
	}
}

func TestRenderContinuationMessage(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	event := events.Event{
		Tick:   960,
		Delta:  480,
		Status: 0xf7,
		Bytes:  []byte{0x83, 0x60, 0xf7, 0x7e, 0x00, 0x09, 0x01, 0xf7},
	}

	message := SysExContinuationMessage{
		SysExEvent: SysExEvent{event},
		Data:       []byte{0x7e, 0x00, 0x09, 0x01, 0xf7},
	}

	w := new(bytes.Buffer)

	message.Render(&ctx, w)

	expected := "   83 60 F7 7E 00 09 01 F7                  tick:960        delta:480        F7 ContinuationMessage 7E 00 09 01 F7"
	if w.String() != expected {
		t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", "SysExContinuationMessage", expected, w.String())
	}
}
