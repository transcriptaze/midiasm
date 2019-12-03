package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestProgramNameRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x08, 0x08, 'P', 'R', 'O', 'G', '-', 'X', 'X', 'X'}},
		0x08,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&ProgramName{metaevent, "PROG-XXX"}, "ProgramName      PROG-XXX"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("ProgramName rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
