package transpose

import (
	"fmt"
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

var Bf3 = midievent.Note{
	Value: 0x3A,
	Name:  "B♭3",
	Alias: "B♭3",
}

var A3x = midievent.Note{
	Value: 0x39,
	Name:  "A3",
	Alias: "A3",
}

var Bf3x = midievent.Note{
	Value: 0x3A,
	Name:  "B♭3",
	Alias: "B♭3",
}

var keySignatureA = metaevent.MakeKeySignature(0, 0, 0, lib.Major, "C Major")
var keySignatureB = metaevent.MakeKeySignature(0, 0, 2, lib.Major, "D♭ Major")
var noteOnA = midievent.MakeNoteOn(0, 0, 1, A3, 76)
var noteOnB = midievent.MakeNoteOn(0, 0, 1, Bf3, 76)
var noteOffA = midievent.MakeNoteOff(0, 0, 1, A3x, 40)
var noteOffB = midievent.MakeNoteOff(0, 0, 1, Bf3x, 40)

var mtrk = midi.MTrk{
	Tag:         "MTrk",
	TrackNumber: 1,
	Events: []*events.Event{
		&events.Event{
			Event: &keySignatureA,
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
			Event: &keySignatureB,
		},
		&events.Event{
			Event: noteOnB,
		},
		&events.Event{
			Event: noteOffB,
		},
		&events.Event{
			Event: metaevent.EndOfTrack{},
		},
	},
	Context: context.NewContext(),
}

func TestTransposeMTrk(t *testing.T) {
	transposed := transpose(mtrk, 1)

	if !reflect.DeepEqual(transposed.Tag, expected.Tag) {
		t.Errorf("Incorrectly transposed MTrk\n   expected:%+v\n   got:     %+v", expected, transposed)
	}

	if !reflect.DeepEqual(transposed.TrackNumber, expected.TrackNumber) {
		t.Errorf("Incorrectly transposed MTrk\n   expected:%+v\n   got:     %+v", expected, transposed)
	}

	if len(transposed.Events) != len(expected.Events) {
		t.Errorf("Incorrectly transposed MTrk\n   expected:%+v\n   got:     %+v", expected, transposed)
	} else {
	}

	for i := range expected.Events {
		p := transposed.Events[i]
		q := expected.Events[i]
		r := mtrk.Events[i]
		fmt.Printf(">>> %#v\n    %#v\n    %#v\n", r.Event, p.Event, q.Event)
		if !reflect.DeepEqual(p, p) {
			t.Errorf("Incorrectly transposed MTrk event\n   expected:%+v\n   got:     %+v", q, p)
		}
	}

	// fmt.Printf(">>> %#v\n    %#v\n", keySignatureA, keySignatureB)
	// fmt.Printf(">>> %#v\n    %#v\n", noteOnA, noteOnB)
}
