package metaevent

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
	{"SequenceNumber",
		&SequenceNumber{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x00, 0x02, 0x12, 0x34}},
				0x00,
			},
			4660},
		"   00 FF 00 02 12 34                        tick:76         delta:12         00 SequenceNumber   4660",
	},

	{"Text",
		&Text{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x08, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}},
				0x01,
			},
			"abcdefgh"},
		"   00 FF 01 08 61 62 63 64 65 66 67 68      tick:76         delta:12         01 Text             abcdefgh",
	},

	{"Copyright",
		&Copyright{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x0b, 'T', 'h', 'e', 'y', ' ', '&', ' ', 'T', 'h', 'e', 'm'}},
				0x02,
			},
			"They & Them"},
		"   00 FF 01 0B 54 68 65 79 20 26 20 54 68\u2026  tick:76         delta:12         02 Copyright        They & Them",
	},

	{"TrackName",
		&TrackName{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x03, 0x0f, 'A', 'c', 'o', 'u', 's', 't', 'i', 'c', ' ', 'G', 'u', 'i', 't', 'a', 'r'}},
				0x03,
			},
			"Acoustic Guitar"},
		"   00 FF 03 0F 41 63 6F 75 73 74 69 63 20\u2026  tick:76         delta:12         03 TrackName        Acoustic Guitar",
	},

	{"InstrumentName",
		&InstrumentName{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x04, 0x06, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72}},
				0x04,
			},
			"Guitar"},
		"   00 FF 04 06 47 75 69 74 61 72            tick:76         delta:12         04 InstrumentName   Guitar",
	},

	{"Lyric",
		&Lyric{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x05, 0x0d, 'L', 'a', 'h', '-', 'l', 'a', '-', 'l', 'a', '-', 'l', 'a', 'h'}},
				0x05,
			},
			"Lah-la-la-lah"},
		"   00 FF 05 0D 4C 61 68 2D 6C 61 2D 6C 61\u2026  tick:76         delta:12         05 Lyric            Lah-la-la-lah",
	},

	{"Marker",
		&Marker{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x06, 0x0b, 'F', 'i', 'r', 's', 't', ' ', 'v', 'e', 'r', 's', 'e'}},
				0x06,
			},
			"First verse"},
		"   00 FF 06 0B 46 69 72 73 74 20 76 65 72\u2026  tick:76         delta:12         06 Marker           First verse",
	},

	{"CuePoint",
		&CuePoint{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x07, 0x0d, 'T', 'h', 'i', 'n', 'g', 's', ' ', 'h', 'a', 'p', 'p', 'e', 'n'}},
				0x07,
			},
			"Things happen"},
		"   00 FF 07 0D 54 68 69 6E 67 73 20 68 61\u2026  tick:76         delta:12         07 CuePoint         Things happen",
	},

	{"ProgramName",
		&ProgramName{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x08, 0x08, 'P', 'R', 'O', 'G', '-', 'X', 'X', 'X'}},
				0x08,
			},
			"PROG-XXX"},
		"   00 FF 08 08 50 52 4F 47 2D 58 58 58      tick:76         delta:12         08 ProgramName      PROG-XXX",
	},

	{"DeviceName",
		&DeviceName{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x09, 0x08, 'D', 'E', 'V', '-', '0', '0', '0', '1'}},
				0x09,
			},
			"DEV-0001"},
		"   00 FF 09 08 44 45 56 2D 30 30 30 31      tick:76         delta:12         09 DeviceName       DEV-0001",
	},

	{"Tempo",
		&Tempo{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20}},
				0x51,
			},
			512438},
		"   00 FF 51 03 07 A1 20                     tick:76         delta:12         51 Tempo            tempo:512438",
	},

	{"EndOfTrack",
		&EndOfTrack{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x2f, 0x00}},
				0x2f,
			},
		},
		"      00 FF 2F 00                           tick:76         delta:12         2F EndOfTrack",
	},

	{"TimeSignature",
		&TimeSignature{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08}},
				0x58,
			},
			4,
			2,
			24,
			8,
		},
		"   00 FF 58 04 04 02 18 08                  tick:76         delta:12         58 TimeSignature    4:4, 24 ticks-per-click, 8/32-per-quarter",
	},

	{"KeySignature",
		&KeySignature{
			MetaEvent{
				event.Event{76, 12, 0xff,
					[]byte{0x00, 0xff, 0x59, 0x02, 0x03, 0x01},
				},
				0x59,
			},
			3,
			1,
		},
		"   00 FF 59 02 03 01                        tick:76         delta:12         59 KeySignature     F\u266f minor",
	},
}

func TestRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if w.String() != v.expected {
			t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.name, v.expected, w.String())
		}
	}
}
