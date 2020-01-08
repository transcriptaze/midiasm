package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestTextRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{"Text", 76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x08, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'}},
		0x01,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&Text{metaevent, "abcdefgh"}, "Text             abcdefgh"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Text rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
