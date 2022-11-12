package midifile

import (
	"bytes"
	// _ "embed"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi"
)

// //go:embed test-files/reference.mid
// var reference []byte

func TestEncodeEmptySMF(t *testing.T) {
	w := bytes.Buffer{}
	smf := midi.SMF{}

	if err := NewEncoder(&w).Encode(smf); err == nil {
		t.Errorf("Expected error, got %v", err)
	}
}

func TestEncodeMThdOnly(t *testing.T) {
	expected := []byte{0x4d, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x00, 0x01, 0xe0}
	w := bytes.Buffer{}
	mthd, _ := midi.NewMThd(1, 0, 480)

	smf := midi.SMF{
		MThd: mthd,
	}

	if err := NewEncoder(&w).Encode(smf); err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(w.Bytes(), expected) {
		t.Errorf("Incorrectly encoded\n   expected:%v\n   got:     %v", hex.Dump(expected), hex.Dump(w.Bytes()))
	}
}

func TestEncodeInvalidMThd(t *testing.T) {
	w := bytes.Buffer{}
	mthd, _ := midi.NewMThd(1, 0, 480)
	track1, _ := midi.NewMTrk()
	track2, _ := midi.NewMTrk()

	smf := midi.SMF{
		MThd:   mthd,
		Tracks: []*midi.MTrk{track1, track2},
	}

	if err := NewEncoder(&w).Encode(smf); err == nil {
		t.Fatalf("Expected error, got %v", err)
	}
}

// func TestEncodeMTrk(t *testing.T) {
// 	expected := []byte{
// 		0x4d, 0x54, 0x72, 0x6b, // MTrk
// 		0x00, 0x00, 0x00, 0x11, // length
// 		0x00, 0xff, 0x03, 0x09, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x20, 0x31, // TrackName
// 		//  0x00, 0xff, 0x51, 0x03, 0x07, 0xa1, 0x20, 0x00, 0xff, 0x2f, 0x00, 0x4d, 0x54
// 		//       0x72, 0x6b, 0x00, 0x00, 0x00, 0x58, 0x00, 0xff, 0x03, 0x0f, 0x41, 0x63, 0x6f, 0x75, 0x73, 0x74
// 		//       0x69, 0x63, 0x20, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72, 0x00, 0xc0, 0x19, 0x00, 0xff, 0x58, 0x04
// 		//       0x04, 0x02, 0x18, 0x08, 0x00, 0x90, 0x30, 0x48, 0x00, 0xff, 0x59, 0x02, 0x00, 0x01, 0x00, 0xb0
// 		//       0x65, 0x00, 0x00, 0xb0, 0x64, 0x00, 0x00, 0xb0, 0x06, 0x06, 0x83, 0x60, 0x80, 0x30, 0x40, 0x00
// 		//       0x90, 0x32, 0x48, 0x83, 0x60, 0x80, 0x32, 0x40, 0x00, 0x90, 0x34, 0x48, 0x83, 0x60, 0x80, 0x34
// 		//       0x40, 0x00, 0x90, 0x35, 0x48, 0x83, 0x60, 0x80, 0x35, 0x40,
// 		0x00, 0xff, 0x2f, 0x00, // EndOfTrack
// 	}
//
// 	mtrk := MTrk{
// 		Tag:    "MTrk",
// 		Length: 24,
// 		Events: []*events.Event{
// 			trackname,
// 			endOfTrack,
// 		},
// 	}
//
// 	bytes, err := mtrk.MarshalBinary()
// 	if err != nil {
// 		t.Fatalf("unexpected error (%v)", err)
// 	}
//
// 	if !reflect.DeepEqual(bytes, expected) {
// 		t.Errorf("incorrectly marshalled\n   expected:%#v\n   got:     %#v", expected, bytes)
// 	}
// }
