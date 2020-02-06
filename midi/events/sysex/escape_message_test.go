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

func TestParseEscapeMessage(t *testing.T) {
	ctx := context.NewContext()
	e := events.Event{
		Status: 0xf7,
		Bytes:  []byte{0x00, 0xf7},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{0x02, 0xf3, 0x01}))

	event, err := Parse(e, r, ctx)
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
