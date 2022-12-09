package events

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
)

func TestEventUnmarshalText(t *testing.T) {
	expected := Event{
		Event: metaevent.MakeSequenceNumber(0, 480, 23, []byte{}...),
	}

	bytes := []byte(`      00 FF 00 02 00 17                     tick:0          delta:480         00 SequenceNumber         23`)
	event := Event{}

	if err := event.UnmarshalText(bytes); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", "SequenceNumber", err)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%#v\n   got:     %#v", "SequenceNumber", expected, event)
	}
}
