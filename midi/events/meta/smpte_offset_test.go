package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestSMPTEOffsetRender(t *testing.T) {
	metaevent := MetaEvent{
		events.Event{"SMPTEOffset", 76, 12, 0xff, []byte{0x00, 0xff, 0x54, 0x05, 9, 2, 5, 28, 13}},
		0x54,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&SMPTEOffset{metaevent, 0x10, 9, 2, 5, 28, 13}, "SMPTEOffset      30fps (drop frame), 09:02:05, 28 frames, 13 fractional frames"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("SMPTEOffset rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
