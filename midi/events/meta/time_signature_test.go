package metaevent

import (
	"bufio"
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestParseTimeSignature(t *testing.T) {
	e := events.Event{
		Tag:    "TimeSignature",
		Tick:   0,
		Delta:  0,
		Status: 0xff,
		Bytes:  []byte{0x00, 0xff},
	}

	r := bufio.NewReader(bytes.NewReader([]byte{4, 3, 3, 24, 8}))

	event, err := NewTimeSignature(&MetaEvent{e, 0x58}, r)
	if err != nil {
		t.Fatalf("TimeSignature parse error: %v", err)
	}

	if event == nil {
		t.Fatalf("TimeSignature parse error - returned %v", event)
	}

	if event.Numerator != 3 {
		t.Errorf("Invalid TimeSignature numerator - expected:%v, got: %v", 3, event.Numerator)
	}

	if event.Denominator != 8 {
		t.Errorf("Invalid TimeSignature denominator - expected:%v, got: %v", 8, event.Denominator)
	}
}

func TestTimeSignatureRender(t *testing.T) {
	metaevent := MetaEvent{
		events.Event{"TimeSignature", 76, 12, 0xff, []byte{0x00, 0xff, 0x58, 0x04, 0x18, 0x08}},
		0x58,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&TimeSignature{metaevent, 1, 1, 24, 8}, "TimeSignature    1/1, 24 ticks per click, 8/32 per quarter"},
		{&TimeSignature{metaevent, 1, 2, 24, 8}, "TimeSignature    1/2, 24 ticks per click, 8/32 per quarter"},
		{&TimeSignature{metaevent, 2, 4, 24, 8}, "TimeSignature    2/4, 24 ticks per click, 8/32 per quarter"},
		{&TimeSignature{metaevent, 3, 4, 24, 8}, "TimeSignature    3/4, 24 ticks per click, 8/32 per quarter"},
		{&TimeSignature{metaevent, 4, 4, 24, 8}, "TimeSignature    4/4, 24 ticks per click, 8/32 per quarter"},
		{&TimeSignature{metaevent, 5, 4, 24, 8}, "TimeSignature    5/4, 24 ticks per click, 8/32 per quarter"},
		{&TimeSignature{metaevent, 6, 8, 24, 8}, "TimeSignature    6/8, 24 ticks per click, 8/32 per quarter"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("TimeSignature rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
