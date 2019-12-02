package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestMarkerRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x06, 0x0b, 'F', 'i', 'r', 's', 't', ' ', 'v', 'e', 'r', 's', 'e'}},
		0x06,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&Marker{metaevent, "First verse"}, "Marker           First verse"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Marker rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
