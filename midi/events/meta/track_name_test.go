package metaevent

import (
	"reflect"
	"testing"
)

func TestTrackNameMarshalBinary(t *testing.T) {
	trackname := TrackName{
		event: event{
			tick:   2400,
			delta:  480,
			Tag:    "TrackName",
			Status: 0xff,
			Type:   0x03,
			bytes:  []byte{},
		},
		Name: "Railroad Traque",
	}

	expected := []byte{
		0xff, 0x03, 0x0f, 0x52, 0x61, 0x69, 0x6c, 0x72,
		0x6f, 0x61, 0x64, 0x20, 0x54, 0x72, 0x61, 0x71,
		0x75, 0x65,
	}

	encoded, err := trackname.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding TrackName (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded TrackName\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTrackNameUnmarshalText(t *testing.T) {
	text := "      00 FF 03 0F 41 63 6F 75 73 74 69 63…  tick:0          delta:0          03 TrackName              Railroad Traque"
	expected := TrackName{
		event: event{
			tick:   0,
			delta:  0,
			Tag:    "TrackName",
			Status: 0xff,
			Type:   0x03,
			bytes:  []byte{},
		},
		Name: "Railroad Traque",
	}

	trackname := TrackName{}

	if err := trackname.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling TrackName (%v)", err)
	}

	if !reflect.DeepEqual(trackname, expected) {
		t.Errorf("incorrectly unmarshalled TrackName\n   expected:%+v\n   got:     %+v", expected, trackname)
	}

}
