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

func TestParseContinuationMessage(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
		Casio: true,
	}

	e := events.Event{
		Tick:   0,
		Delta:  0,
		Status: 0xf7,
		Bytes:  []byte{0x00, 0xf7},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}))

	event, err := Parse(e, r, &ctx)
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
