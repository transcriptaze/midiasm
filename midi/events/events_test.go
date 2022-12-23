package events

import (
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestEventIs(t *testing.T) {
	e := []Event{
		Event{
			Event: metaevent.MakeEndOfTrack(0, 480, bytes[lib.TagEndOfTrack]...),
		},

		Event{
			Event: metaevent.MakeTempo(0, 480, 500000, bytes[lib.TagTempo]...),
		},
	}

	if !Is[metaevent.EndOfTrack](e[0]) {
		t.Errorf("Incorrectly Is'd EndOfTrack event")
	}

	if Is[metaevent.Tempo](e[0]) {
		t.Errorf("Incorrectly Is'd not an EndOfrack event")
	}

	if Is[metaevent.EndOfTrack](e[1]) {
		t.Errorf("Incorrectly Is'd not a Tempo event")
	}

	if !Is[metaevent.Tempo](e[1]) {
		t.Errorf("Incorrectly Is'd Tempo event")
	}
}
