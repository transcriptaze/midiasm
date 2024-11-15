package midifile

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/lib"
)

var SMF0 = []byte{
	0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x60,

	0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x6f,
	0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
	0x00, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27,
	0x00, 0xff, 0x00, 0x02, 0x00, 0x17,
	0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74,
	0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d,
	0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
	0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f,
	0x00, 0xff, 0x59, 0x02, 0x00, 0x01,
	0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e,
	0x00, 0x91, 0x31, 0x48,
	0x00, 0x3c, 0x4c,
	0x00, 0x81, 0x31, 0x64,
	0x00, 0xff, 0x2f, 0x00,
}

var SMF1 = []byte{
	0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60,

	0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x21,
	0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
	0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
	0x00, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27,
	0x00, 0xff, 0x2f, 0x00,

	0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x5f,
	0x00, 0xff, 0x00, 0x02, 0x00, 0x17,
	0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74,
	0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d,
	0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
	0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f,
	0x00, 0xff, 0x59, 0x02, 0x00, 0x01,
	0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e,
	0x00, 0x91, 0x31, 0x48,
	0x00, 0x3c, 0x4c,
	0x00, 0x81, 0x31, 0x64,
	0x00, 0xff, 0x2f, 0x00,
}

var MTHD0 = midi.MThd{
	Tag:      "MThd",
	Length:   6,
	Format:   0,
	Tracks:   1,
	Division: 96,
	PPQN:     96,
	Bytes:    []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x60},
}

var MTHD1 = midi.MThd{
	Tag:      "MThd",
	Length:   6,
	Format:   1,
	Tracks:   2,
	Division: 96,
	PPQN:     96,
	Bytes:    []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60},
}

var MTRK0 = []*midi.MTrk{
	&midi.MTrk{
		Tag:         "MTrk",
		TrackNumber: 0,
		Length:      111,
		Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x6f},
		Events: []*events.Event{
			tempo,
			smpteOffset,
			sequenceNumber,
			text,
			copyright,
			acousticGuitar,
			didgeridoo,
			aMinor,
			motu,
			noteOnCS3,
			noteOnC4,
			noteOffCS3,
			endOfTrack,
		},
	},
}

var MTRK1 = []*midi.MTrk{
	&midi.MTrk{
		Tag:         "MTrk",
		TrackNumber: 0,
		Length:      33,
		Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x21},
		Events: []*events.Event{
			example1,
			tempo,
			smpteOffset,
			endOfTrack,
		},
	},

	&midi.MTrk{
		Tag:         "MTrk",
		TrackNumber: 1,
		Length:      95,
		Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x5f},
		Events: []*events.Event{
			sequenceNumber,
			text,
			copyright,
			acousticGuitar,
			didgeridoo,
			aMinor,
			motu,
			noteOnCS3,
			noteOnC4,
			noteOffCS3,
			endOfTrack,
		},
	},
}

func TestDecodeFormat0(t *testing.T) {
	testDecode(t, SMF0, &MTHD0, MTRK0)
}

func TestDecodeFormat1(t *testing.T) {
	testDecode(t, SMF1, &MTHD1, MTRK1)
}

func testDecode(t *testing.T, b []byte, mthd *midi.MThd, tracks []*midi.MTrk) {
	decoder := NewDecoder()

	smf, err := decoder.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("unexpected error decoding valid MIDI file: %v", err)
	}

	if smf == nil {
		t.Fatalf("decoder returned a 'nil' result for MIDI file")
	}

	if !reflect.DeepEqual(*smf.MThd, *mthd) {
		t.Errorf("MThd incorrectly unmarshaled\n   expected:%v\n   got:     %v", *mthd, *smf.MThd)
	}

	for i, track := range tracks {
		mtrk := smf.Tracks[i]

		if mtrk.Tag != track.Tag {
			t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Tag'\n   expected:%v\n   got:     %v", i, track.Tag, mtrk.Tag)
		}

		if mtrk.TrackNumber != tracks[i].TrackNumber {
			t.Errorf("MTrk[%d]: incorrectly unmarshaled 'TrackNumber'\n   expected:%v\n   got:     %v", i, track.TrackNumber, mtrk.TrackNumber)
		}

		if mtrk.Length != tracks[i].Length {
			t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Length'\n   expected:%v\n   got:     %v", i, track.Length, mtrk.Length)
		}

		if !reflect.DeepEqual(mtrk.Bytes[0:8], tracks[i].Bytes[0:8]) {
			t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Bytes'\n   expected:%v\n   got:     %v", i, track.Bytes[0:8], mtrk.Bytes[0:8])
		}

		if len(mtrk.Events) != len(tracks[i].Events) {
			t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Events'\n   expected:%v\n   got:     %v", i, len(track.Events), len(mtrk.Events))
		} else {
			for j, e := range mtrk.Events {
				if !reflect.DeepEqual(e, tracks[i].Events[j]) {
					t.Errorf("MTrk[%d]: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", i, tracks[i].Events[j], e)
					t.Errorf("\n   expected:%#v\n   got:     %#v", tracks[i].Events[j].Event, e.Event)
				}
			}
		}
	}
}

// TEST EVENTS

var sequenceNumber = &events.Event{
	Event: metaevent.MakeSequenceNumber(0, 0, 23, []byte{0x00, 0xff, 0x00, 0x02, 0x00, 0x17}...),
}

var text = &events.Event{
	Event: metaevent.MakeText(0, 0, "This and That", []byte{0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74}...),
}

var copyright = &events.Event{
	metaevent.MakeCopyright(0, 0, "Them", []byte{0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d}...),
}

var example1 = &events.Event{
	Event: metaevent.MakeTrackName(0, 0, "Example 1", []byte{0x0, 0xff, 0x3, 0x9, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31}...),
}

var acousticGuitar = &events.Event{
	Event: metaevent.MakeTrackName(0, 0, "Acoustic Guitar", []byte{0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72}...),
}

var didgeridoo = &events.Event{
	Event: metaevent.MakeInstrumentName(0, 0, "Didgeridoo", []byte{0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f}...),
}

var aMinor = &events.Event{
	Event: metaevent.MakeKeySignature(0, 0, 0, lib.Minor, []byte{0x00, 0xff, 0x59, 0x02, 0x00, 0x01}...),
}

var motu = &events.Event{
	Event: metaevent.MakeSequencerSpecificEvent(
		0,
		0,
		lib.Manufacturer{
			ID:     []byte{0x00, 0x00, 0x3b},
			Region: "American",
			Name:   "Mark Of The Unicorn (MOTU)"},
		[]byte{0x3a, 0x4c, 0x5e},
		[]byte{0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e}...),
}

var noteOnCS3 = &events.Event{
	Event: midievent.MakeNoteOn(0, 0, 1, midievent.Note{49, "C♯3", "C♯3"}, 72, []byte{0x00, 0x91, 0x31, 0x48}...),
}

var noteOnC4 = &events.Event{
	Event: midievent.MakeNoteOn(0, 0, 1, midievent.Note{60, "C4", "C4"}, 76, []byte{0x00, 0x3c, 0x4c}...),
}

var noteOffCS3 = &events.Event{
	Event: midievent.MakeNoteOff(0, 0, 1, midievent.Note{49, "C♯3", "C♯3"}, 100, []byte{0x00, 0x81, 0x31, 0x64}...),
}

var tempo = &events.Event{
	Event: metaevent.MakeTempo(0, 0, 500000, []byte{0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20}...),
}

var smpteOffset = &events.Event{
	Event: metaevent.MakeSMPTEOffset(0, 0, 13, 45, 59, 25, 7, 39, []byte{0x00, 0xff, 0x54, 0x05, 0x2d, 0x2d, 0x3b, 0x07, 0x27}...),
}

var endOfTrack = &events.Event{
	Event: metaevent.MakeEndOfTrack(0, 0, []byte{0x00, 0xff, 0x2f, 0x00}...),
}
