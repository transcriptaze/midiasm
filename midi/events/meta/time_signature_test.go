package metaevent

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"testing"
)

func TestParseTimeSignature(t *testing.T) {
	e := events.Event{
		Status: 0xff,
		Bytes:  []byte{0x00, 0xff},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{4, 3, 3, 24, 8}))

	event, err := NewTimeSignature(&e, 0x58, r)
	if err != nil {
		t.Fatalf("TimeSignature parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("TimeSignature parse error - returned %v", event)
	}

	if event.Numerator != 3 {
		t.Errorf("Invalid TimeSignature numerator - expected:%v, got: %v", 3, event.Numerator)
	}

	if event.Denominator != 8 {
		t.Errorf("Invalid TimeSignature denominator - expected:%v, got: %v", 8, event.Denominator)
	}
}
