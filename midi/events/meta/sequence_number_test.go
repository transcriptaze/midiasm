package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestSequenceNumberRender(t *testing.T) {
	metaevent := MetaEvent{
		events.Event{"SequenceNumber", 76, 12, 0xff, []byte{0x00, 0xff, 0x00, 0x02, 0x12, 0x34}},
		0x00,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&SequenceNumber{metaevent, 4660}, "SequenceNumber   4660"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("SequenceNumber rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
