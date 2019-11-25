package midievent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"testing"
)

var events = []struct {
	name     string
	event    event.IEvent
	expected string
}{
	{"NoteOff",
		&NoteOff{
			MidiEvent{
				event.Event{1920, 480, 0x80, []byte{0x83, 0x60, 0x80, 0x35, 0x40}},
				0,
			},
			53, 64,
		},
		"   83 60 80 35 40                                            tick:1920       delta:480        80 NoteOff          channel:0 note:53 velocity:64",
	},

	{"NoteOn",
		&NoteOn{
			MidiEvent{
				event.Event{1440, 0, 0x90, []byte{0x00, 0x90, 0x35, 0x48}},
				0,
			},
			53, 72,
		},
		"      00 90 35 48                                            tick:1440       delta:0          90 NoteOn           channel:0 note:53 velocity:72",
	},

	{"PolyphonicPressure",
		&PolyphonicPressure{
			MidiEvent{
				event.Event{1440, 480, 0xa0, []byte{0x00, 0xa0, 0x07}},
				0,
			},
			7,
		},
		"         00 A0 07                                            tick:1440       delta:480        A0 PolyphonicPressure channel:0 pressure:7",
	},

	{"Controller",
		&Controller{
			MidiEvent{
				event.Event{1440, 480, 0xb0, []byte{0x00, 0xb0, 0x06, 0x08}},
				0,
			},
			6, 8,
		},
		"      00 B0 06 08                                            tick:1440       delta:480        B0 Controller       channel:0 controller:6 value:8",
	},

	{"ProgramChange",
		&ProgramChange{
			MidiEvent{
				event.Event{0, 0, 0xc0, []byte{0x00, 0xc0, 0x19}},
				0,
			},
			25,
		},
		"         00 C0 19                                            tick:0          delta:0          C0 ProgramChange    channel:0 program:25",
	},

	{"ChannelPressure",
		&ChannelPressure{
			MidiEvent{
				event.Event{0, 0, 0xd0, []byte{0x00, 0xd0, 0x05}},
				0,
			},
			5,
		},
		"         00 D0 05                                            tick:0          delta:0          D0 ChannelPressure  channel:0 pressure:5",
	},

	{"PitchBend",
		&PitchBend{
			MidiEvent{
				event.Event{0, 0, 0xe0, []byte{0x00, 0xe0, 0x00, 0x08}},
				0,
			},
			8,
		},
		"      00 E0 00 08                                            tick:0          delta:0          E0 PitchBend        channel:0 bend:8",
	},
}

func TestRender(t *testing.T) {
	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if w.String() != v.expected {
			t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.name, v.expected, w.String())
		}
	}
}
