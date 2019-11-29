package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestCopyrightRender(t *testing.T) {
	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x01, 0x02, 0x48, 0x47, 0x46, 0x45, 0x44, 0x43, 0x42, 0x41}},
		0x01,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&Copyright{metaevent, "hgfedcba"}, "Copyright        hgfedcba"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("Copyright rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
