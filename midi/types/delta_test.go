package types

import (
	"testing"
)

func TestDeltaUnmarshalText(t *testing.T) {
	tests := []struct {
		text     string
		expected Delta
	}{
		{"delta:0", 0},
		{"delta:480", 480},
		{"qwerty delta:480", 480},
		{"delta:480 uiop", 480},
		{"qwerty delta:480 uiop", 480},
	}

	for _, test := range tests {
		var delta Delta
		if err := delta.UnmarshalText([]byte(test.text)); err != nil {
			t.Fatalf("Error unmarshalling %v (%v)", test.text, err)
		} else if delta != test.expected {
			t.Errorf("Incorrect delta for %q - expected:%v, got:%v", test.text, test.expected, delta)
		}
	}
}
