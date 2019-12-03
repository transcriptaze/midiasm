package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestDeviceNameRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x09, 0x08, 'D', 'E', 'V', '-', '0', '0', '0', '1'}},
		0x09,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&DeviceName{metaevent, "DEV-0001"}, "DeviceName       DEV-0001"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("DeviceName rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
