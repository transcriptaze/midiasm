package midi

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/events/meta"
	"github.com/twystd/midiasm/midi/types"
	"reflect"
	"testing"
)

func TestUnmarshalSMF(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x02, 0x00, 0x60,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0d,
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x32,
		0x00, 0xff, 0x00, 0x02, 0x00, 0x17,
		0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74,
		0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d,
		0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72,
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
			Length:      13,
			Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0d},
			Events: []events.IEvent{
				&metaevent.TrackName{
					MetaEvent: metaevent.MetaEvent{
						Event: events.Event{
							Tag:    "TrackName",
							Status: 0xff,
							Bytes:  types.Hex{0x0, 0xff, 0x3, 0x9, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31},
						},
						Type: types.MetaEventType(0x03),
					},
					Name: "Example 1",
				},
			},
		},

		MTrk{
			Tag:         "MTrk",
			TrackNumber: 1,
			Length:      50,
			Bytes:       []byte{0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x32},
			Events: []events.IEvent{
				&metaevent.SequenceNumber{
					MetaEvent: metaevent.MetaEvent{
						Event: events.Event{
							Tag:    "SequenceNumber",
							Status: 0xff,
							Bytes:  types.Hex{0x00, 0xff, 0x00, 0x02, 0x00, 0x17},
						},
						Type: types.MetaEventType(0x00),
					},
					SequenceNumber: 23,
				},

				&metaevent.Text{
					MetaEvent: metaevent.MetaEvent{
						Event: events.Event{
							Tag:    "Text",
							Status: 0xff,
							Bytes:  types.Hex{0x00, 0xff, 0x01, 0x0d, 0x54, 0x68, 0x69, 0x73, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x61, 0x74},
						},
						Type: types.MetaEventType(0x01),
					},
					Text: "This and That",
				},

				&metaevent.Copyright{
					MetaEvent: metaevent.MetaEvent{
						Event: events.Event{
							Tag:    "Copyright",
							Status: 0xff,
							Bytes:  types.Hex{0x00, 0xff, 0x02, 0x04, 0x54, 0x68, 0x65, 0x6d},
						},
						Type: types.MetaEventType(0x02),
					},
					Copyright: "Them",
				},

				&metaevent.TrackName{
					MetaEvent: metaevent.MetaEvent{
						Event: events.Event{
							Tag:    "TrackName",
							Status: 0xff,
							Bytes:  types.Hex{0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74, 0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72},
						},
						Type: types.MetaEventType(0x03),
					},
					Name: "Acoustic Guitar",
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
