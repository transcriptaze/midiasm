package midi

import (
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"testing"
)

func TestMTrkPrint(t *testing.T) {
	expected := "4D 54 72 6B 00 00 00 18…                    MTrk 3  length:24"

	mtrk := MTrk{
		Tag:         "MTrk",
		TrackNumber: 3,
		Length:      24,
		Bytes: []byte{
			0x4D, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x18,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Events: []events.IEvent{},
	}

	w := new(bytes.Buffer)

	mtrk.Print(w)

	if w.String() != expected {
		t.Errorf("%s rendered incorrectly\nExpected: '%s'\ngot:      '%s'", "MTrk", expected, w.String())
	}
}
