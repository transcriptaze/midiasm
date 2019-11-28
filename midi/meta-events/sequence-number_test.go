package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestSequenceNumberRender(t *testing.T) {
	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x00, 0x02, 0x12, 0x34}},
		0x00,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&SequenceNumber{metaevent, 4660}, "SequenceNumber   4660"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("SequenceNumber rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
