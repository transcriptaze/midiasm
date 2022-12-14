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
		{
			bytes: []byte{0x83, 0x60, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d},
			expected: Event{
				Event: metaevent.MakeCopyright(0, 480, "Them", []byte{0x83, 0x60, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d}...),
			},
		},
		{
			bytes: []byte{0x83, 0x60, 0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c, 0x72, 0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61, 0x71, 0x75, 0x65},
			expected: Event{
				Event: metaevent.MakeTrackName(0, 480, "Railroad Traque", []byte{0x83, 0x60, 0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c, 0x72, 0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61, 0x71, 0x75, 0x65}...),
			},
		},
		{
			bytes: []byte{0x83, 0x60, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
			expected: Event{
				Event: metaevent.MakeInstrumentName(0, 480, "Didgeridoo", []byte{0x83, 0x60, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f}...),
			},
		},
		{
			bytes: []byte{0x83, 0x60, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61},
			expected: Event{
				Event: metaevent.MakeLyric(0, 480, "La-la-la", []byte{0x83, 0x60, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61}...),
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
