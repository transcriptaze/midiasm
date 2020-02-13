package metaevent

import (
	"testing"
)

func TestParseTimeSignature(t *testing.T) {
	event, err := NewTimeSignature([]byte{3, 3, 24, 8})
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
