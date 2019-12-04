package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestInstrumentNameRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x04, 0x06, 0x47, 0x75, 0x69, 0x74, 0x61, 0x72}},
		0x04,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&InstrumentName{metaevent, "Guitar"}, "InstrumentName   Guitar"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("InstrumentName rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
