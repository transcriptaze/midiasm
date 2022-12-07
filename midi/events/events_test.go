package events

import (
	"reflect"
	"testing"

	"github.com/transcriptaze/midiasm/midi/events/meta"
)

func TestEventUnmarshalJSON(t *testing.T) {
	expected := Event{
		Event: metaevent.MakeSequenceNumber(0, 480, 23, []byte{}...),
	}

	bytes := []byte(`{ "event": { "tag":"SequenceNumber","delta":480,"status":255,"type":0,"sequence-number":23 } }`)
	event := Event{}

	if err := event.UnmarshalJSON(bytes); err != nil {
		t.Fatalf("error unmarshalling %v (%v)", "SequenceNumber", err)
	}

	if !reflect.DeepEqual(event, expected) {
		t.Errorf("incorrectly unmarshalled %v\n   expected:%#v\n   got:     %#v", "SequenceNumber", expected, event)
	}
}
