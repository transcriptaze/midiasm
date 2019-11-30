package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestTimeSignatureRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x58, 0x04, 0x18, 0x08}},
		0x58,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&TimeSignature{metaevent, 1, 1, 24, 8}, "TimeSignature    1:2, 24 ticks-per-click, 8/32-per-quarter"},
		{&TimeSignature{metaevent, 2, 2, 24, 8}, "TimeSignature    2:4, 24 ticks-per-click, 8/32-per-quarter"},
		{&TimeSignature{metaevent, 3, 2, 24, 8}, "TimeSignature    3:4, 24 ticks-per-click, 8/32-per-quarter"},
		{&TimeSignature{metaevent, 4, 2, 24, 8}, "TimeSignature    4:4, 24 ticks-per-click, 8/32-per-quarter"},
		{&TimeSignature{metaevent, 5, 2, 24, 8}, "TimeSignature    5:4, 24 ticks-per-click, 8/32-per-quarter"},
		{&TimeSignature{metaevent, 6, 3, 24, 8}, "TimeSignature    6:8, 24 ticks-per-click, 8/32-per-quarter"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("TimeSignature rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
