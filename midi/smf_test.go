package midi

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/events/midi"
	"github.com/twystd/midiasm/midi/types"
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshalSMF(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60,

		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x21,
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
		0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
		0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27,
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

	mthd := MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   1,
		Tracks:   2,
		Division: 96,
		PPQN:     96,
		Bytes:    []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60},
	}

	tracks := []MTrk{
		MTrk{
			Tag:         "MTrk",
			TrackNumber: 0,
			Length:      33,
			Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x21},
			Events: []*events.EventW{
				&events.EventW{
					Bytes: types.Hex{0x0, 0xff, 0x3, 0x9, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31},
					Event: &metaevent.TrackName{
						Tag:    "TrackName",
						Status: 0xff,
						Type:   types.MetaEventType(0x03),
						Name:   "Example 1",
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20},
					Event: &metaevent.Tempo{
						Tag:    "Tempo",
						Status: 0xff,
						Type:   types.MetaEventType(0x51),
						Tempo:  500000,
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27},
					Event: &metaevent.SMPTEOffset{
						Tag:              "SMPTEOffset",
						Status:           0xff,
						Type:             types.MetaEventType(0x54),
						Hour:             13,
						Minute:           45,
						Second:           59,
						FrameRate:        25,
						Frames:           7,
						FractionalFrames: 39,
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x2f, 0x00},
					Event: &metaevent.EndOfTrack{
						Tag:    "EndOfTrack",
						Status: 0xff,
						Type:   types.MetaEventType(0x2f),
					},
				},
			},
		},

		MTrk{
			Tag:         "MTrk",
			TrackNumber: 1,
			Length:      95,
			Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x5f},
			Events: []*events.EventW{
				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x00, 0x02, 0x00, 0x17},
					Event: &metaevent.SequenceNumber{
						Tag:            "SequenceNumber",
						Status:         0xff,
						Type:           types.MetaEventType(0x00),
						SequenceNumber: 23,
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74},
					Event: &metaevent.Text{
						Tag:    "Text",
						Status: 0xff,
						Type:   types.MetaEventType(0x01),
						Text:   "This and That",
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d},
					Event: &metaevent.Copyright{
						Tag:       "Copyright",
						Status:    0xff,
						Type:      types.MetaEventType(0x02),
						Copyright: "Them",
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72},
					Event: &metaevent.TrackName{
						Tag:    "TrackName",
						Status: 0xff,
						Type:   types.MetaEventType(0x03),
						Name:   "Acoustic Guitar",
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x04, 0x0a, 0x44, 0x69, 0x64, 0x67, 0x65, 0x72, 0x69, 0x64, 0x6f, 0x6f},
					Event: &metaevent.InstrumentName{
						Tag:    "InstrumentName",
						Status: 0xff,
						Type:   types.MetaEventType(0x04),
						Name:   "Didgeridoo",
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x59, 0x02, 0x00, 0x01},
					Event: &metaevent.KeySignature{
						Tag:         "KeySignature",
						Status:      0xff,
						Type:        types.MetaEventType(0x59),
						Accidentals: 0,
						KeyType:     1,
						Key:         "A minor",
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e},
					Event: &metaevent.SequencerSpecificEvent{
						Tag:    "SequencerSpecificEvent",
						Status: 0xff,
						Type:   types.MetaEventType(0x7f),
						Manufacturer: types.Manufacturer{
							ID:     []byte{0x00, 0x00, 0x3b},
							Region: "American",
							Name:   "Mark Of The Unicorn (MOTU)",
						},
						Data: []byte{0x3a, 0x4c, 0x5e},
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0x91, 0x31, 0x48},
					Event: &midievent.NoteOn{
						Tag:     "NoteOn",
						Status:  0x91,
						Channel: types.Channel(0x01),
						Note: midievent.Note{
							Value: 49,
							Name:  "C♯2",
							Alias: "C♯2",
						},
						Velocity: 72,
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0x3c, 0x4c},
					Event: &midievent.NoteOn{
						Tag:     "NoteOn",
						Status:  0x91,
						Channel: types.Channel(0x01),
						Note: midievent.Note{
							Value: 60,
							Name:  "C3",
							Alias: "C3",
						},
						Velocity: 76,
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0x81, 0x31, 0x64},
					Event: &midievent.NoteOff{
						Tag:     "NoteOff",
						Status:  0x81,
						Channel: types.Channel(0x01),
						Note: midievent.Note{
							Value: 49,
							Name:  "C♯2",
							Alias: "C♯2",
						},
						Velocity: 100,
					},
				},

				&events.EventW{
					Bytes: types.Hex{0x00, 0xff, 0x2f, 0x00},
					Event: &metaevent.EndOfTrack{
						Tag:    "EndOfTrack",
						Status: 0xff,
						Type:   types.MetaEventType(0x2f),
					},
				},
			},
		},
	}

	smf := SMF{}
	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling SMF: %v", err)
	}

	if !reflect.DeepEqual(*smf.MThd, mthd) {
		t.Errorf("MThd incorrectly unmarshaled\n   expected:%v\n   got:     %v", mthd, *smf.MThd)
	}

	if len(smf.Tracks) != len(tracks) {
		t.Errorf("MTrk incorrectly unmarshaled 'Tracks'\n   expected:%v\n   got:     %v", len(tracks), len(smf.Tracks))
	} else {
		for i, mtrk := range smf.Tracks {
			if mtrk.Tag != tracks[i].Tag {
				t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Tag'\n   expected:%v\n   got:     %v", i, tracks[i].Tag, mtrk.Tag)
			}

			if mtrk.TrackNumber != tracks[i].TrackNumber {
				t.Errorf("MTrk[%d]: incorrectly unmarshaled 'TrackNumber'\n   expected:%v\n   got:     %v", i, tracks[i].TrackNumber, mtrk.TrackNumber)
			}

			if mtrk.Length != tracks[i].Length {
				t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Length'\n   expected:%v\n   got:     %v", i, tracks[i].Length, mtrk.Length)
			}

			if !reflect.DeepEqual(mtrk.Bytes[0:8], tracks[i].Bytes[0:8]) {
				t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Bytes'\n   expected:%v\n   got:     %v", i, tracks[i].Bytes[0:8], mtrk.Bytes[0:8])
			}

			if len(mtrk.Events) != len(tracks[i].Events) {
				t.Errorf("MTrk[%d]: incorrectly unmarshaled 'Events'\n   expected:%v\n   got:     %v", i, len(tracks[i].Events), len(mtrk.Events))
			} else {
				for j, e := range mtrk.Events {
					if !reflect.DeepEqual(e, tracks[i].Events[j]) {
						t.Errorf("MTrk[%d]: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", i, tracks[i].Events[j], e)
					}
				}
			}
		}
	}
}

func TestUnmarshalSMFWithConf(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x01, 0x00, 0x60,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0a,
		0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e,
	}

	expected := events.EventW{
		Bytes: types.Hex{0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e},
		Event: &metaevent.SequencerSpecificEvent{
			Tag:    "SequencerSpecificEvent",
			Status: 0xff,
			Type:   types.MetaEventType(0x7f),
			Manufacturer: types.Manufacturer{
				ID:     []byte{0x00, 0x00, 0x3b},
				Region: "Borneo",
				Name:   "MOTU",
			},
			Data: []byte{0x3a, 0x4c, 0x5e},
		}}

	conf := `{
  "manufacturers": [
    {
      "id": [ 0, 0, 59 ],
      "region": "Borneo",
      "name": "MOTU"
    }
  ]
}`

	smf := SMF{}

	r := strings.NewReader(conf)
	if err := smf.LoadConfiguration(r); err != nil {
		t.Fatalf("Unexpected error loading configuration (%v)", err)
	}

	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling SMF: %v", err)
	}

	e := smf.Tracks[0].Events[0]
	if !reflect.DeepEqual(e, &expected) {
		t.Errorf("MTrk: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", &expected, e)
	}
}

func TestUnmarshalSMFUnglobalizedConf(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x01, 0x00, 0x60,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0a,
		0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e,
	}

	conf := `{ "manufacturers": [ { "id": [ 0, 0, 59 ], "region": "Borneo", "name": "MOTU" } ] }`

	expected := types.Manufacturer{
		ID:     []byte{0x00, 0x00, 0x3b},
		Region: "Borneo",
		Name:   "MOTU",
	}

	reverted := types.Manufacturer{
		ID:     []byte{0x00, 0x00, 0x3b},
		Region: "American",
		Name:   "Mark Of The Unicorn (MOTU)",
	}

	smf := SMF{}

	r := strings.NewReader(conf)
	if err := smf.LoadConfiguration(r); err != nil {
		t.Fatalf("Unexpected error loading configuration (%v)", err)
	}

	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling SMF: %v", err)
	}

	e := smf.Tracks[0].Events[0].Event.(*metaevent.SequencerSpecificEvent).Manufacturer
	if !reflect.DeepEqual(e, expected) {
		t.Errorf("MTrk: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", expected, e)
	}

	smf = SMF{}

	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling SMF: %v", err)
	}

	e = smf.Tracks[0].Events[0].Event.(*metaevent.SequencerSpecificEvent).Manufacturer
	if !reflect.DeepEqual(e, reverted) {
		t.Errorf("MTrk: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", reverted, e)
	}
}

func TestUnmarshalSMFNoteAlias(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x01, 0x00, 0x60,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x14,
		0x00, 0xff, 0x59, 0x02, 0x06, 0x00,
		0x00, 0x91, 0x31, 0x48,
		0x00, 0xff, 0x59, 0x02, 0xfa, 0x01,
		0x00, 0x81, 0x31, 0x64,
	}

	mtrk := MTrk{
		Tag:         "MTrk",
		TrackNumber: 1,
		Length:      74,
		Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x4a},
		Events: []*events.EventW{
			&events.EventW{
				Bytes: types.Hex{0x00, 0xff, 0x59, 0x02, 0x06, 0x00},
				Event: &metaevent.KeySignature{
					Tag:         "KeySignature",
					Status:      0xff,
					Type:        types.MetaEventType(0x59),
					Accidentals: 6,
					KeyType:     0,
					Key:         "F♯ major",
				},
			},

			&events.EventW{
				Bytes: types.Hex{0x00, 0x91, 0x31, 0x48},
				Event: &midievent.NoteOn{
					Tag:     "NoteOn",
					Status:  0x91,
					Channel: types.Channel(0x01),
					Note: midievent.Note{
						Value: 49,
						Name:  "C♯2",
						Alias: "C♯2",
					},
					Velocity: 72,
				},
			},

			&events.EventW{
				Bytes: types.Hex{0x00, 0xff, 0x59, 0x02, 0xfa, 0x01},
				Event: &metaevent.KeySignature{
					Tag:         "KeySignature",
					Status:      0xff,
					Type:        types.MetaEventType(0x59),
					Accidentals: -6,
					KeyType:     1,
					Key:         "E♭ minor",
				},
			},

			&events.EventW{
				Bytes: types.Hex{0x00, 0x81, 0x31, 0x64},
				Event: &midievent.NoteOff{
					Tag:     "NoteOff",
					Status:  0x81,
					Channel: types.Channel(0x01),
					Note: midievent.Note{
						Value: 49,
						Name:  "C♯2",
						Alias: "D♭2",
					},
					Velocity: 100,
				},
			},
		},
	}

	smf := SMF{}
	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling SMF: %v", err)
	}

	for i, e := range mtrk.Events {
		if !reflect.DeepEqual(e, smf.Tracks[0].Events[i]) {
			t.Errorf("MTrk: incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", e, smf.Tracks[0].Events[i])
		}
	}
}

func TestValidateSMF(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60,

		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x04,
		0x00, 0xff, 0x2f, 0x00,

		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x14,
		0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20,
		0x00, 0xff, 0x54, 0x05, 0x4d, 0x2d, 0x3b, 0x07, 0x27,
		0x00, 0xff, 0x2f, 0x00,
	}

	smf := SMF{}
	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling SMF: %v", err)
	}

	expected := []ValidationError{
		ValidationError(fmt.Errorf("Track 1: unexpected event (Tempo)")),
		ValidationError(fmt.Errorf("Track 1: unexpected event (SMPTEOffset)")),
	}

	errors := smf.Validate()
	if len(errors) != len(expected) {
		t.Errorf("Validation returned %d errors, expected: %v", len(errors), len(expected))
	}

loop:
	for _, e := range expected {
		for _, err := range errors {
			if reflect.DeepEqual(err, e) {
				continue loop
			}
		}
		t.Errorf("Missing expected error: %v", e)
	}

loop2:
	for _, e := range errors {
		for _, err := range expected {
			if reflect.DeepEqual(err, e) {
				continue loop2
			}
		}
		t.Errorf("Unexpected error: %v", e)
	}
}

func TestUnmarshalSMFWithInvalidRunningStatus(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60,

		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x04,
		0x00, 0xff, 0x2f, 0x00,

		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0b,
		0x00, 0x3c, 0x4c,
		0x00, 0x91, 0x31, 0x48,
		0x00, 0xff, 0x2f, 0x00,
	}

	expected := fmt.Errorf("Unrecognised MIDI event: 30")

	smf := SMF{}
	err := smf.UnmarshalBinary(bytes)
	if err == nil {
		t.Fatalf("Expected error unmarshaling SMF - got: %v", nil)
	}

	if !reflect.DeepEqual(err, expected) {
		t.Fatalf("Incorrect error unmarshaling SMF:\nexpected: %+v\n     got: %+v", expected, err)
	}
}
