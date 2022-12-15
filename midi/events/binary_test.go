package events

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/lib"
)

var bytes = map[lib.MetaEventType][]byte{
	lib.TypeSequenceNumber: []byte{0x83, 0x60, 0xff, 0x00, 0x02, 0x00, 0x17},
	lib.TypeText:           []byte{0x83, 0x60, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74},
	lib.TypeCopyright:      []byte{0x83, 0x60, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d},
	lib.TypeTrackName:      []byte{0x83, 0x60, 0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c, 0x72, 0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61, 0x71, 0x75, 0x65},
	lib.TypeInstrumentName: []byte{0x83, 0x60, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
	lib.TypeLyric:          []byte{0x83, 0x60, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61},
	lib.TypeMarker:         []byte{0x83, 0x60, 0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73},
	lib.TypeCuePoint:       []byte{0x83, 0x60, 0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c},
	lib.TypeProgramName:    []byte{0x83, 0x60, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65},
}

func TestEventUnmarshalBinary(t *testing.T) {
	tests := []struct {
		bytes    []byte
		expected Event
	}{
		{
			bytes: bytes[lib.TypeSequenceNumber],
			expected: Event{
				Event: metaevent.MakeSequenceNumber(0, 480, 23, bytes[lib.TypeSequenceNumber]...),
			},
		},
		{
			bytes: bytes[lib.TypeText],
			expected: Event{
				Event: metaevent.MakeText(0, 480, "This and That", bytes[lib.TypeText]...),
			},
		},
		{
			bytes: bytes[lib.TypeCopyright],
			expected: Event{
				Event: metaevent.MakeCopyright(0, 480, "Them", bytes[lib.TypeCopyright]...),
			},
		},
		{
			bytes: bytes[lib.TypeTrackName],
			expected: Event{
				Event: metaevent.MakeTrackName(0, 480, "Railroad Traque", bytes[lib.TypeTrackName]...),
			},
		},
		{
			bytes: bytes[lib.TypeInstrumentName],
			expected: Event{
				Event: metaevent.MakeInstrumentName(0, 480, "Didgeridoo", bytes[lib.TypeInstrumentName]...),
			},
		},
		{
			bytes: bytes[lib.TypeLyric],
			expected: Event{
				Event: metaevent.MakeLyric(0, 480, "La-la-la", bytes[lib.TypeLyric]...),
			},
		},
		{
			bytes: bytes[lib.TypeMarker],
			expected: Event{
				Event: metaevent.MakeMarker(0, 480, "Here Be Dragons", bytes[lib.TypeMarker]...),
			},
		},
		{
			bytes: bytes[lib.TypeCuePoint],
			expected: Event{
				Event: metaevent.MakeCuePoint(0, 480, "More cowbell", bytes[lib.TypeCuePoint]...),
			},
		},
		{
			bytes: bytes[lib.TypeProgramName],
			expected: Event{
				Event: metaevent.MakeProgramName(0, 480, "Escape", bytes[lib.TypeProgramName]...),
			},
		},
	}

	for _, test := range tests {
		event := Event{}

		if err := event.UnmarshalBinary(test.bytes); err != nil {
			t.Fatalf("error unmarshalling %T (%v)", event.Event, err)
		}

		if !reflect.DeepEqual(event, test.expected) {
			t.Errorf("incorrectly unmarshalled %T\n   expected:%#v\n   got:     %#v", event.Event, test.expected, event)
		}
	}
}
