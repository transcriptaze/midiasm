package sysex

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
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

	expected := types.Hex{0x7e, 0x00, 0x09, 0x01, 0xf7}
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

	expected := types.Hex([]byte{0x7e, 0x00, 0x09, 0x01, 0x43})
	if !reflect.DeepEqual(message.Data, expected) {
		t.Errorf("Invalid SysEx single message data - expected:%v, got: %v", expected, message.Data)
	}

	if !ctx.Casio {
		t.Errorf("context Casio flag should be set")
	}
}
