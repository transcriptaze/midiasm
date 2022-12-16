package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalEndOfTrack(t *testing.T) {
	expected := EndOfTrack{
		event: event{
			tick:   2400,
			delta:  480,
			bytes:  []byte{0xff, 0x2f, 0x00},
			tag:    lib.TagEndOfTrack,
			Status: 0xff,
			Type:   0x2f,
		},
	}

	ctx := context.NewContext()
	e := EndOfTrack{}

	err := e.unmarshal(ctx, 2400, 480, 0xff, []byte{}, []byte{0xff, 0x2f, 0x00}...)

	if err != nil {
		t.Fatalf("error encoding EndOfTrack (%v)", err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrect EndOfTrack\n   expected:%+v\n   got:     %+v", expected, e)
	}
}

func TestEndOfTrackMarshalBinary(t *testing.T) {
	evt := EndOfTrack{
		event: event{
			tick:   2400,
			delta:  0,
			tag:    lib.TagEndOfTrack,
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

func TestEndOfTrackUnmarshalBinary(t *testing.T) {
	expected := EndOfTrack{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagEndOfTrack,
			Status: 0xff,
			Type:   lib.TypeEndOfTrack,
			bytes:  []byte{0x83, 0x60, 0xff, 0x2f, 0x00},
		},
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x2f, 0x00}

	e := EndOfTrack{}

	if err := e.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error encoding %v (%v)", lib.TagEndOfTrack, err)
	}

	if !reflect.DeepEqual(e, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", lib.TagEndOfTrack, expected, e)
	}
}

func TestEndOfTrackUnmarshalText(t *testing.T) {
	text := "      00 FF 2F 00                           tick:0          delta:480        2F EndOfTrack"
	expected := EndOfTrack{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagEndOfTrack,
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
			tag:    lib.TagEndOfTrack,
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
			tag:    lib.TagEndOfTrack,
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
