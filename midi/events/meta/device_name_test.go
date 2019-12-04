package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestDeviceNameRender(t *testing.T) {
	ctx := context.Context{
		Scale: context.Sharps,
	}

	metaevent := MetaEvent{
		events.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x09, 0x08, 'D', 'E', 'V', '-', '0', '0', '0', '1'}},
		0x09,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&DeviceName{metaevent, "DEV-0001"}, "DeviceName       DEV-0001"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("DeviceName rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
