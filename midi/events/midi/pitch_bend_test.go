package midievent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

func TestParsePitchBend(t *testing.T) {
	expected := PitchBend{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xe7, 0x00, 0x08},

			tag:     types.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
		},
		Bend: 8,
	}

	ctx := context.NewContext()
	r := IO.TestReader([]byte{0x00, 0xe7}, []byte{0x00, 0x08})

	event, err := Parse(2400, 480, r, 0xe7, ctx)
	if err != nil {
		t.Fatalf("Unexpected PitchBend event parse error: %v", err)
	} else if event == nil {
		t.Fatalf("Unexpected PitchBend event parse error - returned %v", event)
	}

	v, ok := event.(*PitchBend)
	if !ok {
		t.Fatalf("PitchBend %v type error - returned %T", event, v)
	} else if !reflect.DeepEqual(event, &expected) {
		t.Errorf("Invalid PitchBend event\n  expected:%#v\n  got:     %#v", &expected, event)
	}
}

func TestPitchBendMarshalBinary(t *testing.T) {
	evt := PitchBend{
		event: event{
			tick:  2400,
			delta: 480,
			bytes: []byte{0x00, 0xe7, 0x00, 0x08},
			tag:   types.TagPitchBend,

			Status:  0xe7,
			Channel: 7,
		},
		Bend: 8,
	}

	expected := []byte{0xe7, 0x00, 0x08}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding PitchBend (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded PitchBend\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestPitchBendUnmarshalText(t *testing.T) {
	text := "      81 70 E7 00 08                           tick:240        delta:240        E7 PitchBend              channel:7  bend:8"
	expected := PitchBend{
		event: event{
			tick:    0,
			delta:   240,
			tag:     types.TagPitchBend,
			Status:  0xe7,
			Channel: 7,
			bytes:   []byte{},
		},
		Bend: 8,
	}

	evt := PitchBend{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling PitchBend (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled PitchBend\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}