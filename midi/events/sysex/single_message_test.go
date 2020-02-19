package sysex

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"reflect"
	"testing"
)

func TestParseSingleMessage(t *testing.T) {
	ctx := context.NewContext()
	r := bufio.NewReader(bytes.NewReader([]byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0xf7}))

	event, err := Parse(reader{r}, 0xf0, ctx)
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

	manufacturer := types.Manufacturer{
		ID:     []byte{0x7e},
		Region: "Special Purpose",
		Name:   "Non-RealTime Extensions",
	}
	if !reflect.DeepEqual(message.Manufacturer, manufacturer) {
		t.Errorf("Invalid SysEx single message manufacturer - expected:%v, got: %v", manufacturer, message.Manufacturer)
	}

	data := types.Hex{0x00, 0x09, 0x01}
	if !reflect.DeepEqual(message.Data, data) {
		t.Errorf("Invalid SysEx single message data - expected:%v, got: %v", data, message.Data)
	}

	if ctx.Casio {
		t.Errorf("context Casio flag should not be set")
	}
}

func TestParseSingleMessageWithoutTerminatingF7(t *testing.T) {
	ctx := context.NewContext()
	r := bufio.NewReader(bytes.NewReader([]byte{0x05, 0x7e, 0x00, 0x09, 0x01, 0x43}))

	event, err := Parse(reader{r}, 0xf0, ctx)
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

	manufacturer := types.Manufacturer{
		ID:     []byte{0x7e},
		Region: "Special Purpose",
		Name:   "Non-RealTime Extensions",
	}
	if !reflect.DeepEqual(message.Manufacturer, manufacturer) {
		t.Errorf("Invalid SysEx single message manufacturer - expected:%v, got: %v", manufacturer, message.Manufacturer)
	}

	data := types.Hex([]byte{0x00, 0x09, 0x01, 0x43})
	if !reflect.DeepEqual(message.Data, data) {
		t.Errorf("Invalid SysEx single message data - expected:%v, got: %v", data, message.Data)
	}

	if !ctx.Casio {
		t.Errorf("context Casio flag should be set")
	}
}
