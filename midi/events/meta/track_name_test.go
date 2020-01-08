package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestTrackNameRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{"TrackName", 76, 12, 0xff, []byte{0x00, 0xff, 0x03, 0x0f, 'A', 'c', 'o', 'u', 's', 't', 'i', 'c', ' ', 'G', 'u', 'i', 't', 'a', 'r'}},
		0x03,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&TrackName{metaevent, "Acoustic Guitar"}, "TrackName        Acoustic Guitar"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("TrackName rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
