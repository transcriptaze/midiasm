package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestMarkerRender(t *testing.T) {
	metaevent := MetaEvent{
		events.Event{"Marker", 76, 12, 0xff, []byte{0x00, 0xff, 0x06, 0x0b, 'F', 'i', 'r', 's', 't', ' ', 'v', 'e', 'r', 's', 'e'}},
		0x06,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&Marker{metaevent, "First verse"}, "Marker           First verse"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Marker rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
