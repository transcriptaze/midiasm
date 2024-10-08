package notes

import (
	"bytes"
	_ "embed"
	"reflect"
	"testing"
	"time"

	"github.com/transcriptaze/midiasm/encoding/midi"
	"github.com/transcriptaze/midiasm/midi"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/events/meta"
	"github.com/transcriptaze/midiasm/midi/events/midi"
)

var smf = midi.SMF{
	MThd: &midi.MThd{
		Format:        1,
		Tracks:        2,
		PPQN:          480,
		Division:      480,
		SMPTETimeCode: false,
		SubFrames:     0,
		DropFrame:     false,
	},

	Tracks: []*midi.MTrk{
		&midi.MTrk{
			Events: []*events.Event{
				&events.Event{
					Event: metaevent.MakeTempo(0, 0, 500000),
				},
			},
		},
		&midi.MTrk{
			Events: []*events.Event{
				&events.Event{
					Event: midievent.MakeNoteOn(0, 0, 0, midievent.Note{
						Value: 48,
						Name:  "C3",
						Alias: "C3",
					}, 72),
				},
				&events.Event{
					Event: midievent.MakeNoteOff(480, 480, 0, midievent.Note{
						Value: 48,
						Name:  "C3",
						Alias: "C3",
					}, 64),
				},
				&events.Event{
					Event: midievent.MakeNoteOn(480, 0, 0, midievent.Note{
						Value: 50,
						Name:  "D3",
						Alias: "D3",
					}, 72),
				},
				&events.Event{
					Event: midievent.MakeNoteOff(960, 480, 0, midievent.Note{
						Value: 50,
						Name:  "D3",
						Alias: "D3",
					}, 64),
				},
				&events.Event{
					Event: midievent.MakeNoteOn(960, 0, 0, midievent.Note{
						Value: 52,
						Name:  "E3",
						Alias: "E3",
					}, 72),
				},
				&events.Event{
					Event: midievent.MakeNoteOff(1440, 480, 0, midievent.Note{
						Value: 52,
						Name:  "E3",
						Alias: "E3",
					}, 64),
				},
				&events.Event{
					Event: midievent.MakeNoteOn(1440, 0, 0, midievent.Note{
						Value: 53,
						Name:  "F3",
						Alias: "F3",
					}, 72),
				},
				&events.Event{
					Event: midievent.MakeNoteOff(1920, 480, 0, midievent.Note{
						Value: 53,
						Name:  "F3",
						Alias: "F3",
					}, 64),
				},
				&events.Event{
					Event: midievent.MakeNoteOn(1920, 0, 0, midievent.Note{
						Value: 55,
						Name:  "G3",
						Alias: "G3",
					}, 72),
				},
				&events.Event{ // NoteOn with zero velocity
					Event: midievent.MakeNoteOn(2400, 480, 0, midievent.Note{
						Value: 55,
						Name:  "G3",
						Alias: "G3",
					}, 0),
				},
			},
		},
	},
}

func TestExtractNotes(t *testing.T) {
	expected := []Note{
		Note{
			Channel:       0,
			Note:          48,
			FormattedNote: "C3",
			Velocity:      72,
			StartTick:     0,
			EndTick:       480,
			Start:         0 * time.Millisecond,
			End:           500 * time.Millisecond,
			Duration:      500 * time.Millisecond,
		},
		Note{
			Channel:       0,
			Note:          50,
			FormattedNote: "D3",
			Velocity:      72,
			StartTick:     480,
			EndTick:       960,
			Start:         500 * time.Millisecond,
			End:           1000 * time.Millisecond,
			Duration:      500 * time.Millisecond,
		},
		Note{
			Channel:       0,
			Note:          52,
			FormattedNote: "E3",
			Velocity:      72,
			StartTick:     960,
			EndTick:       1440,
			Start:         1000 * time.Millisecond,
			End:           1500 * time.Millisecond,
			Duration:      500 * time.Millisecond,
		},
		Note{
			Channel:       0,
			Note:          53,
			FormattedNote: "F3",
			Velocity:      72,
			StartTick:     1440,
			EndTick:       1920,
			Start:         1500 * time.Millisecond,
			End:           2000 * time.Millisecond,
			Duration:      500 * time.Millisecond,
		},
		Note{
			Channel:       0,
			Note:          55,
			FormattedNote: "G3",
			Velocity:      72,
			StartTick:     1920,
			EndTick:       2400,
			Start:         2000 * time.Millisecond,
			End:           2500 * time.Millisecond,
			Duration:      500 * time.Millisecond,
		},
	}

	notes, err := extract(&smf, 0)
	if err != nil {
		t.Fatalf("Error extracting notes from SMF (%v)", err)
	}

	if len(notes) != len(expected) {
		t.Errorf("Incorrectly extracted notes - expected:%v notes, got: %v notes", len(expected), len(notes))
	} else {
		for i := range expected {
			p := expected[i]
			q := notes[i]
			if !reflect.DeepEqual(p, q) {
				t.Errorf("Incorrectly extracted note %v\n   expected:%v\n   got:     %v", i+1, p, q)
			}
		}
	}

}

//go:embed test-files/notes.mid
var testfile []byte

//go:embed test-files/notes.txt
var reference []byte

func TestNotesWithTempoChanges(t *testing.T) {
	smf, _ := midifile.NewDecoder().Decode(bytes.NewReader(testfile))

	notes, err := extract(smf, 0)
	if err != nil {
		t.Fatalf("Error extracting notes from SMF (%v)", err)
	}

	var b bytes.Buffer
	if err := print(notes, &b); err != nil {
		t.Fatalf("Error extracting notes from SMF (%v)", err)
	}

	if !reflect.DeepEqual(b.Bytes(), reference) {
		t.Errorf("incorrectly extracted notes\nexpected:\n%+v\ngot:\n%+v", string(reference), b.String())
	}
}

//go:embed test-files/notes-with-missing-note-offs.mid
var testfile2 []byte

//go:embed test-files/notes-with-missing-note-offs.txt
var reference2 []byte

func TestExtractNotesWithMissingNoteOff(t *testing.T) {
	smf, _ := midifile.NewDecoder().Decode(bytes.NewReader(testfile2))

	notes, err := extract(smf, 0)
	if err != nil {
		t.Fatalf("Error extracting notes from SMF (%v)", err)
	}

	var b bytes.Buffer
	if err := print(notes, &b); err != nil {
		t.Fatalf("Error extracting notes from SMF (%v)", err)
	}

	if !reflect.DeepEqual(b.Bytes(), reference2) {
		t.Errorf("incorrectly extracted notes\nexpected:\n%+v\ngot:\n%+v", string(reference2), b.String())
	}
}
