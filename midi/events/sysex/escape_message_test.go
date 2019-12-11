package sysex

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"reflect"
	"testing"
)

func TestParseEscapeMessage(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
		Casio: false,
	}

	e := events.Event{
		Tick:   0,
		Delta:  0,
		Status: 0xf7,
		Bytes:  []byte{0x00, 0xf7},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x02, 0xf3, 0x01}))

	event, err := Parse(e, r, &ctx)
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

	expected := []byte{0xf3, 0x01}
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx escape message data - expected:%v, got: %v", expected, message.Data)
	}
}

func TestRenderEscapeMessage(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
		Casio: false,
	}

	event := events.Event{
		Tick:   960,
		Delta:  480,
		Status: 0xf7,
		Bytes:  []byte{0x83, 0x60, 0xf7, 0x02, 0xf3, 0x01},
	}

	message := SysExEscapeMessage{
		SysExEvent: SysExEvent{event},
		Data:       []byte{0xf3, 0x01},
	}

	w := new(bytes.Buffer)

	message.Render(&ctx, w)

	expected := "   83 60 F7 02 F3 01                        tick:960        delta:480        F7 EscapeMessage    F3 01"
	if w.String() != expected {
		t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", "SysExEscapeMessage", expected, w.String())
	}
}
