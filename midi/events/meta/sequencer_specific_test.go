package metaevent

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

func TestUnmarshalSequencerSpecificEvent(t *testing.T) {
	expected := SequencerSpecificEvent{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
			bytes:  []byte{0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x00, 0x00, 0x3b},
			Region: "American",
			Name:   "Mark Of The Unicorn (MOTU)",
		},
		Data: []byte{0x3a, 0x4c, 0x5e},
	}

	ctx := context.NewContext()

	evt, err := UnmarshalSequencerSpecificEvent(ctx, 2400, 480, []byte{0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e}, []byte{0x00, 0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e}...)
	if err != nil {
		t.Fatalf("error unmarshalling SequencerSpecificEvent (%v)", err)
	}

	if !reflect.DeepEqual(*evt, expected) {
		t.Errorf("incorrect SequencerSpecificEvent\n   expected:%+v\n   got:     %+v", expected, *evt)
	}
}

func TestSequencerSpecificEventMarshalBinary(t *testing.T) {
	evt := SequencerSpecificEvent{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
			bytes:  []byte{},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x00, 0x00, 0x3b},
			Region: "American",
			Name:   "Mark Of The Unicorn (MOTU)",
		},
		Data: []byte{0x3a, 0x4c, 0x5e},
	}

	expected := []byte{0xff, 0x7f, 0x06, 0x00, 0x00, 0x3b, 0x3a, 0x4c, 0x5e}

	encoded, err := evt.MarshalBinary()
	if err != nil {
		t.Fatalf("error encoding SequencerSpecific (%v)", err)
	}

	if !reflect.DeepEqual(encoded, expected) {
		t.Errorf("incorrectly encoded SequencerSpecific\n   expected:%+v\n   got:     %+v", expected, encoded)
	}
}

func TestTextUnmarshalSequencerSpecific(t *testing.T) {
	text := "      00 FF 7F 06 00 00 3B 3A 4C 5E         tick:0          delta:480        7F SequencerSpecificEvent Mark Of The Unicorn (MOTU), 3A 4C 5E"
	expected := SequencerSpecificEvent{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
			bytes:  []byte{},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x00, 0x00, 0x3b},
			Region: "American",
			Name:   "Mark Of The Unicorn (MOTU)",
		},
		Data: []byte{0x3a, 0x4c, 0x5e},
	}

	evt := SequencerSpecificEvent{}

	if err := evt.UnmarshalText([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling SequencerSpecificEvent (%v)", err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled SequencerSpecificEvent\n   expected:%+v\n   got:     %+v", expected, evt)
	}

}

func TestSequencerSpecificEventMarshalJSON(t *testing.T) {
	e := SequencerSpecificEvent{
		event: event{
			tick:   2400,
			delta:  480,
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
			bytes:  []byte{},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x00, 0x00, 0x3b},
			Region: "American",
			Name:   "Mark Of The Unicorn (MOTU)",
		},
		Data: []byte{0x3a, 0x4c, 0x5e},
	}

	expected := `{"tag":"SequencerSpecificEvent","delta":480,"status":255,"type":127,"manufacturer":{"id":[0,0,59],"region":"American","name":"Mark Of The Unicorn (MOTU)"},"data":[58,76,94]}`

	testMarshalJSON(t, lib.TagSequencerSpecificEvent, e, expected)
}

func TestSequencerSpecificEventNameUnmarshalJSON(t *testing.T) {
	tag := lib.TagSequencerSpecificEvent
	text := `{"tag":"SequencerSpecificEvent","delta":480,"status":255,"type":127,"manufacturer":{"id":[0,0,59],"region":"American","name":"Mark Of The Unicorn (MOTU)"},"data":[58,76,94]}`
	expected := SequencerSpecificEvent{
		event: event{
			tick:   0,
			delta:  480,
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
			bytes:  []byte{},
		},
		Manufacturer: lib.Manufacturer{
			ID:     []byte{0x00, 0x00, 0x3b},
			Region: "American",
			Name:   "Mark Of The Unicorn (MOTU)",
		},
		Data: []byte{0x3a, 0x4c, 0x5e},
	}

	evt := SequencerSpecificEvent{}

	if err := evt.UnmarshalJSON([]byte(text)); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", tag, err)
	}

	if !reflect.DeepEqual(evt, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%+v\n   got:     %+v", tag, expected, evt)
	}
}
