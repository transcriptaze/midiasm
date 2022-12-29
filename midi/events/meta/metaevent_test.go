package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type TestMetaEvent interface {
	SequenceNumber |
		Text |
		Copyright |
		TrackName |
		InstrumentName |
		Lyric |
		Marker |
		CuePoint |
		ProgramName |
		DeviceName |
		MIDIChannelPrefix |
		MIDIPort |
		KeySignature |
		SequencerSpecificEvent

	MarshalJSON() ([]byte, error)
}

func TestParseCMajorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{0x83, 0x60, 0xff, 0x59, 0x02, 0x00, 0x00},
		},
		Accidentals: 0,
		KeyType:     lib.Major,
		Key:         "C major",
	}

	event, err := Parse(2400, []byte{0x83, 0x60, 0xff, 0x59, 0x02, 0x00, 0x00}...)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	if ks, ok := event.(KeySignature); !ok {
		t.Fatalf("KeySignature event parse error - returned %T", event)
	} else if !reflect.DeepEqual(ks, expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", expected, ks)
	}
}

func TestParseCMinorKeySignature(t *testing.T) {
	expected := KeySignature{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagKeySignature,
			Status: 0xff,
			Type:   0x59,
			bytes:  []byte{0x83, 0x60, 0xff, 0x59, 0x02, 0xfd, 0x01},
		},
		Accidentals: -3,
		KeyType:     lib.Minor,
		Key:         "C minor",
	}

	event, err := Parse(2400, []byte{0x83, 0x60, 0xff, 0x59, 0x02, 0xfd, 0x01}...)
	if err != nil {
		t.Fatalf("Unexpected KeySignature event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected KeySignature event parse error - returned %v", event)
	}

	if e, ok := event.(KeySignature); !ok {
		t.Fatalf("KeySignature event parse error - returned %T", e)
	} else if !reflect.DeepEqual(e, expected) {
		t.Errorf("Invalid KeySignature event\n  expected:%#v\n  got:     %#v", expected, e)
	}
}

func testMarshalJSON[E TestMetaEvent](t *testing.T, tag lib.Tag, e E, expected string) {
	encoded, err := e.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding %v (%v)", tag, err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded %v\n   expected:%+v\n   got:     %+v", tag, expected, string(encoded))
	}
}
