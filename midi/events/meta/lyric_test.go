package metaevent

import (
	"reflect"
	"testing"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

func TestUnmarshalLyric(t *testing.T) {
	expected := Lyric{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   0x05,
			bytes:  []byte{0x00, 0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61},
		},
		Lyric: "La-la-la",
	}

	evt, err := UnmarshalLyric(2400, 480, []byte("La-la-la"))
	if err != nil {
		t.Fatalf("error encoding Lyric (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect Lyric\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestLyricMarshalBinary(t *testing.T) {
	evt := Lyric{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   0x05,
			bytes:  []byte{},
		},
		Lyric: "La-la-la",
	}

	expected := []byte{0xff, 0x05, 0x08, 0x4c, 0x61, 0x2d, 0x6c, 0x61, 0x2d, 0x6c, 0x61}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding Lyric (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded Lyric\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalLyric(t *testing.T) {
	text := "      00 FF 05 08 4C 61 2D 6C 61 2D 6C 61   tick:0          delta:480        05 Lyric                  La-la-la"
	expected := Lyric{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   0x05,
			bytes:  []byte{},
		},
		Lyric: "La-la-la",
	}

	evt := Lyric{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling Lyric (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled Lyric\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestLyricMarshalJSON(t *testing.T) {
	e := Lyric{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   lib.TypeLyric,
			bytes:  []byte{},
		},
		Lyric: "La-la-la",
	}

	expected := `{"tag":"Lyric","delta":480,"status":255,"type":5,"lyric":"La-la-la"}`

	testMarshalJSON(t, lib.TagLyric, e, expected)
}

func TestLyricNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagLyric
	text := `{"tag":"Lyric","delta":480,"status":255,"type":5,"lyric":"La-la-la"}`
	expected := Lyric{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   lib.TypeLyric,
			bytes:  []byte{},
		},
		Lyric: "La-la-la",
	}

	evt := Lyric{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
