package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestSMPTEOffsetRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x54, 0x05, 9, 2, 5, 28, 13}},
		0x54,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&SMPTEOffset{metaevent, 0x10, 9, 2, 5, 28, 13}, "SMPTEOffset      30fps (drop frame), 09:02:05, 28 frames, 13 fractional frames"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("SMPTEOffset rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
