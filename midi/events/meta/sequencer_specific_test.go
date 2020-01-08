package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestSequencerSpecificEventRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{"SequencerSpecificEvent", 76, 12, 0xff, []byte{0x00, 0xff, 0x7f, 0x06, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab}},
		0x7f,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&SequencerSpecificEvent{metaevent, []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab}}, "SequencerSpecificEvent 01 23 45 67 89 AB"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("SequencerSpecificEvent rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
