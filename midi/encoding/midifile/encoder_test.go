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
