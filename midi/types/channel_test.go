package types

import (
	"testing"
)

func TestParseChannel(t *testing.T) {
	tests := []struct {
		text     string
		expected Channel
	}{
		{"0", 0},
		{"12", 12},
		{"13", 13},
		{"15", 15},
	}

	for _, test := range tests {
		if channel, err := ParseChannel(test.text); err != nil {
			t.Fatalf("Error parsing %v (%v)", test.text, err)
		} else if channel != test.expected {
			t.Errorf("Incorrect channel for %q - expected:%v, got:%v", test.text, test.expected, channel)
		}
	}
}

func TestParseInvalidChannel(t *testing.T) {
	tests := []struct {
		text     string
		expected Channel
	}{
		{"-1", 0},
		{"16", 16},
	}

	for _, test := range tests {
		if channel, err := ParseChannel(test.text); err == nil {
			t.Errorf("Expected error parsing invalid channel %v, got (%v, %v)", test.text, channel, err)
		}
	}
}

func TestChannelUnmarshalText(t *testing.T) {
	tests := []struct {
		text     string
		expected Channel
	}{
		{"channel:0", 0},
		{"qwerty channel:12", 12},
		{"channel:13 uiop", 13},
		{"qwerty channel:15 uiop", 15},
	}

	for _, test := range tests {
		var channel Channel
		if err := channel.UnmarshalText([]byte(test.text)); err != nil {
			t.Fatalf("Error unmarshalling %v (%v)", test.text, err)
		} else if channel != test.expected {
			t.Errorf("Incorrect channel for %q - expected:%v, got:%v", test.text, test.expected, channel)
		}
	}
}
