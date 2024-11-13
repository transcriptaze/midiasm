package events

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/events/sysex"
	"github.com/transcriptaze/midiasm/midi/lib"
)

var bytes = map[lib.Tag][]byte{
	lib.TagSequenceNumber:         []byte{0x83, 0x60, 0xff, 0x00, 0x02, 0x00, 0x17},
	lib.TagText:                   []byte{0x83, 0x60, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74},
	lib.TagCopyright:              []byte{0x83, 0x60, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d},
	lib.TagTrackName:              []byte{0x83, 0x60, 0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c, 0x72, 0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61, 0x71, 0x75, 0x65},
	lib.TagInstrumentName:         []byte{0x83, 0x60, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
	lib.TagLyric:                  []byte{0x83, 0x60, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61},
	lib.TagMarker:                 []byte{0x83, 0x60, 0xff, 0x06, 0x0f, 0x48, 0x65, 0x72, 0x65, 0x20, 0x42, 0x65, 0x20, 0x44, 0x72, 0x61, 0x67, 0x6f, 0x6e, 0x73},
	lib.TagCuePoint:               []byte{0x83, 0x60, 0xff, 0x07, 0x0c, 0x4d, 0x6f, 0x72, 0x65, 0x20, 0x63, 0x6f, 0x77, 0x62, 0x65, 0x6c, 0x6c},
	lib.TagProgramName:            []byte{0x83, 0x60, 0xff, 0x08, 0x06, 0x45, 0x73, 0x63, 0x61, 0x70, 0x65},
	lib.TagDeviceName:             []byte{0x83, 0x60, 0xff, 0x09, 0x08, 0x54, 0x68, 0x65, 0x54, 0x68, 0x69, 0x6e, 0x67},
	lib.TagMIDIChannelPrefix:      []byte{0x83, 0x60, 0xff, 0x20, 0x01, 0x0d},
	lib.TagMIDIPort:               []byte{0x83, 0x60, 0xff, 0x21, 0x01, 0x70},
	lib.TagTempo:                  []byte{0x83, 0x60, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20},
	lib.TagTimeSignature:          []byte{0x83, 0x60, 0xff, 0x58, 0x04, 0x03, 0x02, 0x18, 0x08},
	lib.TagKeySignature:           []byte{0x83, 0x60, 0xff, 0x59, 0x02, 0xfd, 0x01},
	lib.TagSMPTEOffset:            []byte{0x83, 0x60, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27},
	lib.TagEndOfTrack:             []byte{0x83, 0x60, 0xff, 0x2f, 0x00},
	lib.TagSequencerSpecificEvent: []byte{0x83, 0x60, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e},

	lib.TagNoteOff:            []byte{0x83, 0x60, 0x87, 0x31, 0x48},
	lib.TagNoteOn:             []byte{0x83, 0x60, 0x97, 0x31, 0x48},
	lib.TagPolyphonicPressure: []byte{0x83, 0x60, 0xa7, 0x64},
	lib.TagController:         []byte{0x83, 0x60, 0xb7, 0x54, 0x1d},
	lib.TagProgramChange:      []byte{0x83, 0x60, 0xc7, 0x19},
	lib.TagChannelPressure:    []byte{0x83, 0x60, 0xd7, 0x64},
	lib.TagPitchBend:          []byte{0x83, 0x60, 0xe7, 0x00, 0x08},

	lib.TagSysExMessage:      []byte{0x83, 0x60, 0xf0, 0x04, 0x7e, 0x00, 0x09, 0x01},
	lib.TagSysExContinuation: []byte{0x83, 0x60, 0xf7, 0x04, 0x7e, 0x00, 0x09, 0x01},
}

func TestEventUnmarshalBinary(t *testing.T) {
	tests := []struct {
		bytes    []byte
		expected Event
	}{
		{
			bytes: bytes[lib.TagSequenceNumber],
			expected: Event{
				Event: metaevent.MakeSequenceNumber(0, 480, 23, bytes[lib.TagSequenceNumber]...),
			},
		},
		{
			bytes: bytes[lib.TagText],
			expected: Event{
				Event: metaevent.MakeText(0, 480, "This and That", bytes[lib.TagText]...),
			},
		},
		{
			bytes: bytes[lib.TagCopyright],
			expected: Event{
				Event: metaevent.MakeCopyright(0, 480, "Them", bytes[lib.TagCopyright]...),
			},
		},
		{
			bytes: bytes[lib.TagTrackName],
			expected: Event{
				Event: metaevent.MakeTrackName(0, 480, "Railroad Traque", bytes[lib.TagTrackName]...),
			},
		},
		{
			bytes: bytes[lib.TagInstrumentName],
			expected: Event{
				Event: metaevent.MakeInstrumentName(0, 480, "Didgeridoo", bytes[lib.TagInstrumentName]...),
			},
		},
		{
			bytes: bytes[lib.TagLyric],
			expected: Event{
				Event: metaevent.MakeLyric(0, 480, "La-la-la", bytes[lib.TagLyric]...),
			},
		},
		{
			bytes: bytes[lib.TagMarker],
			expected: Event{
				Event: metaevent.MakeMarker(0, 480, "Here Be Dragons", bytes[lib.TagMarker]...),
			},
		},
		{
			bytes: bytes[lib.TagCuePoint],
			expected: Event{
				Event: metaevent.MakeCuePoint(0, 480, "More cowbell", bytes[lib.TagCuePoint]...),
			},
		},
		{
			bytes: bytes[lib.TagProgramName],
			expected: Event{
				Event: metaevent.MakeProgramName(0, 480, "Escape", bytes[lib.TagProgramName]...),
			},
		},
		{
			bytes: bytes[lib.TagDeviceName],
			expected: Event{
				Event: metaevent.MakeDeviceName(0, 480, "TheThing", bytes[lib.TagDeviceName]...),
			},
		},
		{
			bytes: bytes[lib.TagMIDIChannelPrefix],
			expected: Event{
				Event: metaevent.MakeMIDIChannelPrefix(0, 480, 13, bytes[lib.TagMIDIChannelPrefix]...),
			},
		},
		{
			bytes: bytes[lib.TagMIDIPort],
			expected: Event{
				Event: metaevent.MakeMIDIPort(0, 480, 112, bytes[lib.TagMIDIPort]...),
			},
		},
		{
			bytes: bytes[lib.TagTempo],
			expected: Event{
				Event: metaevent.MakeTempo(0, 480, 500000, bytes[lib.TagTempo]...),
			},
		},
		{
			bytes: bytes[lib.TagTimeSignature],
			expected: Event{
				Event: metaevent.MakeTimeSignature(0, 480, 3, 4, 24, 8, bytes[lib.TagTimeSignature]...),
			},
		},
		{
			bytes: bytes[lib.TagKeySignature],
			expected: Event{
				Event: metaevent.MakeKeySignature(0, 480, -3, 1, bytes[lib.TagKeySignature]...),
			},
		},
		{
			bytes: bytes[lib.TagSMPTEOffset],
			expected: Event{
				Event: metaevent.MakeSMPTEOffset(0, 480, 13, 45, 59, 25, 7, 39, bytes[lib.TagSMPTEOffset]...),
			},
		},
		{
			bytes: bytes[lib.TagEndOfTrack],
			expected: Event{
				Event: metaevent.MakeEndOfTrack(0, 480, bytes[lib.TagEndOfTrack]...),
			},
		},
		{
			bytes: bytes[lib.TagSequencerSpecificEvent],
			expected: Event{
				Event: metaevent.MakeSequencerSpecificEvent(0, 480, lib.Manufacturer{
					ID:     []byte{0x00, 0x00, 0x3b},
					Region: "American",
					Name:   "Mark Of The Unicorn (MOTU)",
				}, []byte{0x3a, 0x4c, 0x5e}, bytes[lib.TagSequencerSpecificEvent]...),
			},
		},
		{
			bytes: bytes[lib.TagNoteOff],
			expected: Event{
				Event: midievent.MakeNoteOff(0, 480, 7, midievent.Note{
					Value: 49,
					Name:  "C♯3",
					Alias: "C♯3",
				}, 72, bytes[lib.TagNoteOff]...),
			},
		},
		{
			bytes: bytes[lib.TagNoteOn],
			expected: Event{
				Event: midievent.MakeNoteOn(0, 480, 7, midievent.Note{
					Value: 49,
					Name:  "C♯3",
					Alias: "C♯3",
				}, 72, bytes[lib.TagNoteOn]...),
			},
		},
		{
			bytes: bytes[lib.TagPolyphonicPressure],
			expected: Event{
				Event: midievent.MakePolyphonicPressure(0, 480, 7, 100, bytes[lib.TagPolyphonicPressure]...),
			},
		},
		{
			bytes: bytes[lib.TagController],
			expected: Event{
				Event: midievent.MakeController(0, 480, 7, lib.Controller{84, "Portamento Control"}, 29, bytes[lib.TagController]...),
			},
		},
		{
			bytes: bytes[lib.TagProgramChange],
			expected: Event{
				Event: midievent.MakeProgramChange(0, 480, 7, 0, 25, bytes[lib.TagProgramChange]...),
			},
		},
		{
			bytes: bytes[lib.TagChannelPressure],
			expected: Event{
				Event: midievent.MakeChannelPressure(0, 480, 7, 100, bytes[lib.TagChannelPressure]...),
			},
		},
		{
			bytes: bytes[lib.TagPitchBend],
			expected: Event{
				Event: midievent.MakePitchBend(0, 480, 7, 8, bytes[lib.TagPitchBend]...),
			},
		},
		{
			bytes: bytes[lib.TagSysExMessage],
			expected: Event{
				Event: sysex.MakeSysExMessage(0, 480, lib.Manufacturer{
					ID:     []byte{0x7e},
					Region: "Special Purpose",
					Name:   "Non-RealTime Extensions",
				}, lib.Hex{0x00, 0x09, 0x01}, bytes[lib.TagSysExMessage]...),
			},
		},
		{
			bytes: bytes[lib.TagSysExContinuation],
			expected: Event{
				Event: sysex.MakeSysExContinuationMessage(0, 480, lib.Hex{0x7e, 0x00, 0x09, 0x01}, bytes[lib.TagSysExContinuation]...),
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
