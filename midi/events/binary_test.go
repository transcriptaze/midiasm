package events

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
)

func TestEventUnmarshalBinary(t *testing.T) {
	tests := []struct {
		bytes    []byte
		expected Event
	}{
		{
			bytes: []byte{0x83, 0x60, 0xff, 0x00, 0x02, 0x00, 0x17},
			expected: Event{
				Event: metaevent.MakeSequenceNumber(0, 480, 23, []byte{0x83, 0x60, 0xff, 0x0, 0x2, 0x0, 0x17}...),
			},
		},

		{
			bytes: []byte{0x83, 0x60, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74},
			expected: Event{
				Event: metaevent.MakeText(0, 480, "This and That", []byte{0x83, 0x60, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74}...),
			},
		},
	}

	for _, test := range tests {
		event := Event{}

		if err := event.UnmarshalBinary(test.bytes); err != nil {
			t.Fatalf("error unmarshalling %v (%v)", "SequenceNumber", err)
		}

		if !reflect.DeepEqual(event, test.expected) {
			t.Errorf("incorrectly unmarshalled %v\n   expected:%#v\n   got:     %#v", "SequenceNumber", test.expected, event)
		}
	}
}
