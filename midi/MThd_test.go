package midi

import (
	"reflect"
	"testing"
)

func TestMThdMarshal(t *testing.T) {
	expected := []byte{
		0x4d, 0x54, 0x68, 0x64, // MThd
		0x00, 0x00, 0x00, 0x06, // length
		0x00, 0x01, // format
		0x00, 0x01, // tracks
		0x01, 0xE0, // division
	}

	mthd := MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   1,
		Tracks:   1,
		Division: 480,
	}

	bytes, err := mthd.MarshalBinary()
	if err != nil {
		t.Fatalf("unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(bytes, expected) {
		t.Errorf("incorrectly marshalled\n   expected:%#v\n   got:     %#v", expected, bytes)
	}
}

func TestUnmarshalMThdFormat0(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0x01, 0xe0,
	}

	expected := MThd{
		Tag:      "MThd",
		Length:   6,
		Format:   0,
		Tracks:   1,
		Division: 480,
		PPQN:     480,
		// SMPTETimeCode bool    // TODO make getter/TextUnmarshal
		// SubFrames     uint16  // TODO make getter/TextUnmarshal
		// FPS           uint8   // TODO make getter/TextUnmarshal
		// DropFrame     bool    // TODO make getter/TextUnmarshal
		Bytes: bytes,
	}

	mthd := MThd{}

	if err := mthd.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("unexpected error unmarshaling MThd: %v", err)
	}

	if !reflect.DeepEqual(mthd, expected) {
		t.Errorf("incorrectly unmarshaled MThd chunk\n   expected:%#v\n   got:     %#v", expected, mthd)
	}
}

func TestUnmarshalMThdSMPTE24FPS(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0xe8, 0x50,
	}

	expected := MThd{
		Tag:           "MThd",
		Length:        6,
		Format:        0,
		Tracks:        1,
		Division:      0xe850,
		SMPTETimeCode: true,
		FPS:           24,
		SubFrames:     80,
		Bytes:         bytes,
	}

	mthd := MThd{}

	if err := mthd.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("unexpected error unmarshaling MThd: %v", err)
	}

	if !reflect.DeepEqual(mthd, expected) {
		t.Errorf("incorrectly unmarshaled MThd chunk\n   expected:%#v\n   got:     %#v", expected, mthd)
	}
}

func TestUnmarshalMThdSMPTE25FPS(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0xe7, 0x28,
	}

	expected := MThd{
		Tag:           "MThd",
		Length:        6,
		Format:        0,
		Tracks:        1,
		Division:      0xe728,
		SMPTETimeCode: true,
		FPS:           25,
		SubFrames:     40,
		Bytes:         bytes,
	}

	mthd := MThd{}

	if err := mthd.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("unexpected error unmarshaling MThd: %v", err)
	}

	if !reflect.DeepEqual(mthd, expected) {
		t.Errorf("incorrectly unmarshaled MThd chunk\n   expected:%#v\n   got:     %#v", expected, mthd)
	}
}

func TestUnmarshalMThdSMPTE29FPS(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0xe3, 0x50,
	}

	expected := MThd{
		Tag:           "MThd",
		Length:        6,
		Format:        0,
		Tracks:        1,
		Division:      0xe350,
		SMPTETimeCode: true,
		FPS:           29,
		DropFrame:     true,
		SubFrames:     80,
		Bytes:         bytes,
	}

	mthd := MThd{}

	if err := mthd.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("unexpected error unmarshaling MThd: %v", err)
	}

	if !reflect.DeepEqual(mthd, expected) {
		t.Errorf("incorrectly unmarshaled MThd chunk\n   expected:%#v\n   got:     %#v", expected, mthd)
	}
}

func TestUnmarshalMThdSMPTE30FPS(t *testing.T) {
	bytes := []byte{
		0x4d, 0x54, 0x68, 0x64,
		0x00, 0x00, 0x00, 0x06,
		0x00, 0x00,
		0x00, 0x01,
		0xe2, 0x50,
	}

	expected := MThd{
		Tag:           "MThd",
		Length:        6,
		Format:        0,
		Tracks:        1,
		Division:      0xe250,
		SMPTETimeCode: true,
		FPS:           30,
		SubFrames:     80,
		// DropFrame     bool    // TODO make getter/TextUnmarshal
		Bytes: bytes,
	}

	mthd := MThd{}

	if err := mthd.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("unexpected error unmarshaling MThd: %v", err)
	}

	if !reflect.DeepEqual(mthd, expected) {
		t.Errorf("incorrectly unmarshaled MThd chunk\n   expected:%#v\n   got:     %#v", expected, mthd)
	}
}

// func TestMThdUnmarshalSMTPE(t *testing.T) {
// 	bytes := []byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x11, 0xe7, 0x28}
// 	expected := MThd{
// 		Tag:           "MThd",
// 		Length:        6,
// 		Format:        1,
// 		Tracks:        17,
// 		Division:      0xe728,
// 		PPQN:          0,
// 		SMPTETimeCode: true,
// 		SubFrames:     40,
// 		FPS:           25,
// 		Bytes:         []byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x11, 0xe7, 0x28},
// 	}
//
// 	mthd := MThd{}
// 	if err := mthd.UnmarshalBinary(bytes); err != nil {
// 		t.Fatalf("Unexpected error unmarshaling MThd: %v", err)
// 	}
//
// 	if !reflect.DeepEqual(mthd, expected) {
// 		t.Errorf("MThd incorrectly unmarshaled\n   expected:%+v\n   got:     %+v", expected, mthd)
// 	}
// }

// func TestMThdUnmarshalInvalidBytes(t *testing.T) {
// 	mthd := MThd{}
// 	bytes := [][]byte{
// 		[]byte{0x4D, 0x54, 0x72, 0x6b, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x11, 0x00, 0x60},
// 		[]byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x05, 0x00, 0x01, 0x00, 0x11, 0x00, 0x60},
// 		[]byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x03, 0x00, 0x11, 0x00, 0x60},
// 		[]byte{0x4D, 0x54, 0x68, 0x64, 0x00, 0x00, 0x00, 0x06, 0x00, 0x01, 0x00, 0x11, 0x00},
// 	}
//
// 	for _, b := range bytes {
// 		if err := mthd.UnmarshalBinary(b); err == nil {
// 			t.Fatalf("Expected error unmarshaling MThd: got %v", err)
// 		}
// 	}
// }
