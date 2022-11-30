package transpose

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
	"github.com/transcriptaze/midiasm/midi/lib"
)

var A3 = midievent.Note{
	Value: 0x39,
	Name:  "A3",
	Alias: "A3",
}

var Asharp3 = midievent.Note{
	Value: 0x3A,
	Name:  "A♯3",
	Alias: "A♯3",
}

var Bflat3 = midievent.Note{
	Value: 0x3A,
	Name:  "B♭3",
	Alias: "B♭3",
}

var keySignatureA = metaevent.MakeKeySignature(0, 0, 0, lib.Major)  // "C major"
var keySignatureB = metaevent.MakeKeySignature(0, 0, -5, lib.Major) // "D♭ major"
var noteOnA = midievent.MakeNoteOn(0, 0, 1, A3, 76)
var noteOnAsharp = midievent.MakeNoteOn(0, 0, 1, Asharp3, 76)
var noteOnBflat = midievent.MakeNoteOn(0, 0, 1, Bflat3, 76)
var noteOffA = midievent.MakeNoteOff(0, 0, 1, A3, 40)
var noteOffAsharp = midievent.MakeNoteOff(0, 0, 1, Asharp3, 40)
var noteOffBflat = midievent.MakeNoteOff(0, 0, 1, Bflat3, 40)

var mtrk = midi.MTrk{
	Tag:         "MTrk",
	TrackNumber: 1,
	Events: []*events.Event{
		&events.Event{
			Event: keySignatureA,
		},
		&events.Event{
			Event: noteOnA,
		},
		&events.Event{
			Event: noteOffA,
		},
		&events.Event{
			Event: metaevent.EndOfTrack{},
		},
	},
	Context: context.NewContext(),
}

var expected = midi.MTrk{
	Tag:         "MTrk",
	TrackNumber: 1,
	Events: []*events.Event{
		&events.Event{
			Event: keySignatureB,
		},
		&events.Event{
			Event: noteOnBflat,
		},
		&events.Event{
			Event: noteOffBflat,
		},
		&events.Event{
			Event: metaevent.EndOfTrack{},
		},
	},
	Context: context.NewContext(),
}

func TestTransposeMTrk(t *testing.T) {
	transposed := transpose(mtrk, 1)
	if transposed == nil {
		t.Fatalf("Invalid transposed MTrk (%v)", transposed)
	}

	if !reflect.DeepEqual(transposed.Tag, expected.Tag) {
		t.Errorf("Incorrectly transposed MTrk\n   expected:%+v\n   got:     %+v", expected, *transposed)
	}

	if !reflect.DeepEqual(transposed.TrackNumber, expected.TrackNumber) {
		t.Errorf("Incorrectly transposed MTrk\n   expected:%+v\n   got:     %+v", expected, *transposed)
	}

	if len(transposed.Events) != len(expected.Events) {
		t.Errorf("Incorrectly transposed MTrk\n   expected:%+v\n   got:     %+v", expected, *transposed)
	} else {
	}

	for i := range expected.Events {
		p := transposed.Events[i]
		q := expected.Events[i]
		if !reflect.DeepEqual(p, q) {
			t.Errorf("Incorrectly transposed MTrk event\n   expected:%+v\n   got:     %+v", q, p)
		}
	}
}
