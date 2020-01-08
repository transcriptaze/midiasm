package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestMIDIChannelPrefixRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{"MIDIChannelPrefix", 76, 12, 0xff, []byte{0x00, 0xff, 0x20, 0x01, 13}},
		0x20,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&MIDIChannelPrefix{metaevent, 13}, "MIDIChannelPrefix 13"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("MIDIChannelPrefix rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
