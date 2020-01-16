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
		0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x01, 0x00, 0x60,
		0x4d, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x0d,
		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31,
	}

	mthd := MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   1,
		Tracks:   1,
		Division: 96,
		PPQN:     96,
		Bytes:    []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x01, 0x00, 0x60},
	}

	track0 := MTrk{
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
	}

	smf := SMF{}
	if err := smf.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("Unexpected error unmarshaling MThd: %v", err)
	}

	if !reflect.DeepEqual(*smf.MThd, mthd) {
		t.Errorf("MThd incorrectly unmarshaled\n   expected:%v\n   got:     %v", mthd, *smf.MThd)
	}

	if smf.Tracks[0].Tag != track0.Tag {
		t.Errorf("MTrk[0] incorrectly unmarshaled 'Tag'\n   expected:%v\n   got:     %v", track0.Tag, smf.Tracks[0].Tag)
	}

	if smf.Tracks[0].TrackNumber != track0.TrackNumber {
		t.Errorf("MTrk[0] incorrectly unmarshaled 'TrackNumber'\n   expected:%v\n   got:     %v", track0.TrackNumber, smf.Tracks[0].TrackNumber)
	}

	if smf.Tracks[0].Length != track0.Length {
		t.Errorf("MTrk[0] incorrectly unmarshaled 'Length'\n   expected:%v\n   got:     %v", track0.Length, smf.Tracks[0].Length)
	}

	if !reflect.DeepEqual(smf.Tracks[0].Bytes[0:8], track0.Bytes[0:8]) {
		t.Errorf("MTrk[0] incorrectly unmarshaled 'Bytes'\n   expected:%v\n   got:     %v", track0.Bytes[0:8], smf.Tracks[0].Bytes[0:8])
	}

	if len(smf.Tracks[0].Events) != len(track0.Events) {
		t.Errorf("MTrk[0] incorrectly unmarshaled 'Events'\n   expected:%v\n   got:     %v", len(track0.Events), len(smf.Tracks[0].Events))
	} else {
		for i, e := range smf.Tracks[0].Events {
			if !reflect.DeepEqual(e, track0.Events[i]) {
				t.Errorf("MTrk[0] incorrectly unmarshaled event\n   expected:%#v\n   got:     %#v", track0.Events[i], e)
			}
		}
	}
}
