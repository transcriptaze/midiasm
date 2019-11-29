package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestTextRender(t *testing.T) {
	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x08, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48}},
		0x01,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&Text{metaevent, "abcdefgh"}, "Text             abcdefgh"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Text rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
