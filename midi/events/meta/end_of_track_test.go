package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalEndOfTrack(t *testing.T) {
	expected := EndOfTrack{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  nil,
			tag:    types.TagEndOfTrack,
			Status: 0xff,
			Type:   0x2f,
		},
	}

	evt, err := UnmarshalEndOfTrack(2400, 480, []byte{}...)
	if err != nil {
		t.Fatalf("error encoding EndOfTrack (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect EndOfTrack\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestEndOfTrackMarshalBinary(t *testing.T) {
	evt := EndOfTrack{
		event: event{
			tick:   2400,
			delta:  0,
			tag:    types.TagEndOfTrack,
			Status: 0xff,
			Type:   0x2f,
			bytes:  []byte{},
		},
	}

	expected := []byte{0xff, 0x2f, 0x00}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding EndOfTrack (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded EndOfTrack\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestEndOfTrackUnmarshalText(t *testing.T) {
	text := "      00 FF 2F 00                           tick:0          delta:480        2F EndOfTrack"
	expected := EndOfTrack{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagEndOfTrack,
			Status: 0xff,
			Type:   0x2f,
			bytes:  []byte{},
		},
	}

	evt := EndOfTrack{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling EndOfTrack (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled EndOfTrack\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestEndOfTrackMarshalJSON(t *testing.T) {
	evt := EndOfTrack{
		event: event{
			tick:   2400,
			delta:  0,
			tag:    types.TagEndOfTrack,
			Status: 0xff,
			Type:   0x2f,
			bytes:  []byte{},
		},
	}

	expected := `{"tag":"EndOfTrack","delta":0,"status":255,"type":47}`

	encoded, err := evt.MarshalJSON()
	if err != nil {
		t.Fatalf("error encoding EndOfTrack (%v)", err)
	}

	if string(encoded) != expected {
		t.Errorf("incorrectly encoded EndOfTrack\n   expected:%v\n   got:     %v", expected, string(encoded))
	}
}

func TestEndOfTrackUnmarshalJSON(t *testing.T) {
	text := `{"tag":"EndOfTrack","delta":480,"status":255,"type":47}`
	expected := EndOfTrack{
		event: event{
			tick:   0,
			delta:  480,
			tag:    types.TagEndOfTrack,
			Status: 0xff,
			Type:   0x2f,
			bytes:  []byte{},
		},
	}

	evt := EndOfTrack{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling EndOfTrack (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled EndOfTrack\n   expected:%+v\n   got:     %+v", expected, evt)
	}
}
