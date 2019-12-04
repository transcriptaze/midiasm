package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestSequencerSpecificEventRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x7f, 0x06, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab}},
		0x7f,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&SequencerSpecificEvent{metaevent, []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}}, "SequencerSpecificEvent 01 23 45 67 89 AB"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("SequencerSpecificEvent rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
