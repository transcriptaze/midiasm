package events

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
)

func TestEventUnmarshalBinary(t *testing.T) {
	expected := Event{
		Event: metaevent.MakeSequenceNumber(0, 480, 23, []byte{0x83, 0x60, 0xff, 0x0, 0x2, 0x0, 0x17}...),
	}

	bytes := []byte{0x83, 0x60, 0xff, 0x00, 0x02, 0x00, 0x17}
	event := Event{}

	if err := event.UnmarshalBinary(bytes); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", "SequenceNumber", err)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%#v\n   got:     %#v", "SequenceNumber", expected, event)
	}
}
