package metaevent

import (
	"bufio"
	"bytes"
	"testing"
)

func TestParseTimeSignature(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{4, 3, 3, 24, 8}))

	event, err := NewTimeSignature(r)
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
