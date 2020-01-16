package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestCuePointRender(t *testing.T) {
	metaevent := MetaEvent{
		events.Event{"CuePoint", 76, 12, 0xff, []byte{0x00, 0xff, 0x07, 0x0d, 'T', 'h', 'i', 'n', 'g', 's', ' ', 'h', 'a', 'p', 'p', 'e', 'n'}},
		0x07,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&CuePoint{metaevent, "Things happen"}, "CuePoint         Things happen"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("CuePoint rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
