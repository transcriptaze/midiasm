package lib

import (
	"testing"
)

func TestParseDelta(t *testing.T) {
	tests := []struct {
		text     string
		expected Delta
	}{
		{"0", 0},
		{"480", 480},
	}

	for _, test := range tests {
		if delta, err := ParseDelta(test.text); err != nil {
			t.Fatalf("Error parsing %v (%v)", test.text, err)
		} else if delta != test.expected {
			t.Errorf("Incorrect delta for %q - expected:%v, got:%v", test.text, test.expected, delta)
		}
	}
}

func TestParseInvalidDelta(t *testing.T) {
	tests := []struct {
		text     string
		expected Channel
	}{
		{"-1", 0},
	}

	for _, test := range tests {
		if delta, err := ParseDelta(test.text); err == nil {
			t.Errorf("Expected error parsing invalid delta %v, got (%v, %v)", test.text, delta, err)
		}
	}
}
