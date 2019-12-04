package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestLyricRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x05, 0x0d, 'L', 'a', 'h', '-', 'l', 'a', '-', 'l', 'a', '-', 'l', 'a', 'h'}},
		0x05,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&Lyric{metaevent, "Lah-la-la-lah"}, "Lyric            Lah-la-la-lah"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Lyric rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
