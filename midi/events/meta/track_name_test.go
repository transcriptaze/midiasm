package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalTrackName(t *testing.T) {
	expected := TrackName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
			bytes: []byte{
				0x00, 0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c,
				0x72, 0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61,
				0x71, 0x75, 0x65},
		},
		Name: "Railroad Traque",
	}

	ctx := context.NewContext()
	e := TrackName{}

	err := e.unmarshal(ctx, 2400, 480, 0xff, []byte("Railroad Traque"), []byte{
		0x00, 0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c,
		0x72, 0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61,
		0x71, 0x75, 0x65}...)

	if err != nil {
		t.Fatalf("error encoding TrackName (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect TrackName\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestTrackNameMarshalBinary(t *testing.T) {
	evt := TrackName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
			bytes:  []byte{},
		},
		Name: "Railroad Traque",
	}

	expected := []byte{
		0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c, 0x72,
		0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61, 0x71,
		0x75, 0x65,
	}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding TrackName (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded TrackName\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTrackNameUnmarshalText(t *testing.T) {
	text := "      00 FF 03 0F 41 63 6F 75 73 74 69 63…  tick:0          delta:480        03 TrackName              Railroad Traque"
	expected := TrackName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
			bytes:  []byte{},
		},
		Name: "Railroad Traque",
	}

	evt := TrackName{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling TrackName (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled TrackName\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestTrackNameMarshalJSON(t *testing.T) {
	e := TrackName{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
			bytes:  []byte{},
		},
		Name: "Railroad Traque",
	}

	expected := `{"tag":"TrackName","delta":480,"status":255,"type":3,"name":"Railroad Traque"}`

	testMarshalJSON(t, lib.TagTrackName, e, expected)
}

func TestTrackNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagTrackName
	text := `{"tag":"TrackName","delta":480,"status":255,"type":3,"name":"Railroad Traque"}`
	expected := TrackName{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagTrackName,
			Status: 0xff,
			Type:   lib.TypeTrackName,
			bytes:  []byte{},
		},
		Name: "Railroad Traque",
	}

	evt := TrackName{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
