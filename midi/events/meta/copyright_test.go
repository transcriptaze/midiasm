package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestCopyrightRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x0b, 'T', 'h', 'e', 'y', ' ', '&', ' ', 'T', 'h', 'e', 'm'}},
		0x01,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&Copyright{metaevent, "They & Them"}, "Copyright        They & Them"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Copyright rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
