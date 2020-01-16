package sysex

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"reflect"
	"testing"
)

func TestParseSingleMessage(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
		Casio: false,
	}

	e := events.Event{
		Tick:   0,
		Delta:  0,
		Status: 0xf0,
		Bytes:  []byte{0x00, 0xf0},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}))

	event, err := Parse(e, r, &ctx)
	if err != nil {
		t.Fatalf("Unexpected SysEx single message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx single message parse error - returned %v", event)
	}

	message, ok := event.(*SysExSingleMessage)
	if !ok {
		t.Fatalf("SysEx single message parse error - returned %T", event)
	}

	expected := []byte{0x7e, 0x00, 0x09, 0x01, 0xf7}
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx single message data - expected:%v, got: %v", expected, message.Data)
	}

	if ctx.Casio {
		t.Errorf("context Casio flag should not be set")
	}
}

func TestParseSingleMessageWithoutTerminatingF7(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
		Casio: false,
	}

	e := events.Event{
		Tick:   0,
		Delta:  0,
		Status: 0xf0,
		Bytes:  []byte{0x00, 0xf0},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0x43}))

	event, err := Parse(e, r, &ctx)
	if err != nil {
		t.Fatalf("Unexpected SysEx single message parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("Unexpected SysEx single message parse error - returned %v", event)
	}

	message, ok := event.(*SysExSingleMessage)
	if !ok {
		t.Fatalf("SysEx single message parse error - returned %T", event)
	}

	expected := []byte{0x7e, 0x00, 0x09, 0x01, 0x43}
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx single message data - expected:%v, got: %v", expected, message.Data)
	}

	if !ctx.Casio {
		t.Errorf("context Casio flag should be set")
	}
}

func TestRenderSingleMessage(t *testing.T) {
	event := events.Event{
		Tick:   960,
		Delta:  480,
		Status: 0xf0,
		Bytes:  []byte{0x83, 0x60, 0xf0, 0x7e, 0x00, 0x09, 0x01, 0xf7},
	}

	message := SysExSingleMessage{
		SysExEvent: SysExEvent{event},
		Data:       []byte{0x7e, 0x00, 0x09, 0x01, 0xf7},
	}

	w := new(bytes.Buffer)

	message.Render(w)

	expected := "   83 60 F0 7E 00 09 01 F7                  tick:960        delta:480        F0 SingleMessage    7E 00 09 01 F7"
	if w.String() != expected {
		t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", "SysExSingleMessage", expected, w.String())
	}
}
