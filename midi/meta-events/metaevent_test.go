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
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x08, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48}},
				0x01,
			},
			"abcdefgh"},
		"   00 FF 01 08 41 42 43 44 45 46 47 48      tick:76         delta:12         01 Text             abcdefgh",
	},

	{"TrackName",
		&TrackName{
			MetaEvent{
				event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72}},
				0x03,
			},
			"Acoustic Guitar"},
		"   00 FF 03 0F 41 63 6F 75 73 74 69 63 20\u2026  tick:76         delta:12         03 TrackName        Acoustic Guitar",
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
	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if w.String() != v.expected {
			t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.name, v.expected, w.String())
		}
	}
}
