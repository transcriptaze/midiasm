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
	// r := bufio.NewReader(bytes.NewReader([]byte{0x02, 0xf3, 0x01}))
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
