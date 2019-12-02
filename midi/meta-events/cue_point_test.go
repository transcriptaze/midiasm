package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestCuePointRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x07, 0x0d, 'T', 'h', 'i', 'n', 'g', 's', ' ', 'h', 'a', 'p', 'p', 'e', 'n'}},
		0x07,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&CuePoint{metaevent, "Things happen"}, "CuePoint         Things happen"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("CuePoint rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
