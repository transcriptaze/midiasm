package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/events"
	"strings"
	"testing"
)

func TestKeySignatureRender(t *testing.T) {
	metaevent := MetaEvent{
		events.Event{"KeySignature", 76, 12, 0xff, []byte{0x00, 0xff, 0x59, 0x02, 0x00, 0x00}},
		0x59,
	}

	var eventlist = []struct {
		event    events.IEvent
		expected string
	}{
		{&KeySignature{metaevent, 0, 0, "C major"}, "KeySignature     C major"},
		{&KeySignature{metaevent, 1, 0, "G major"}, "KeySignature     G major"},
		{&KeySignature{metaevent, 2, 0, "D major"}, "KeySignature     D major"},
		{&KeySignature{metaevent, 3, 0, "A major"}, "KeySignature     A major"},
		{&KeySignature{metaevent, 4, 0, "E major"}, "KeySignature     E major"},
		{&KeySignature{metaevent, 5, 0, "B major"}, "KeySignature     B major"},
		{&KeySignature{metaevent, 6, 0, "F\u266f major"}, "KeySignature     F\u266f major"},
		{&KeySignature{metaevent, -1, 0, "F major"}, "KeySignature     F major"},
		{&KeySignature{metaevent, -2, 0, "B\u266d major"}, "KeySignature     B\u266d major"},
		{&KeySignature{metaevent, -3, 0, "E\u266d major"}, "KeySignature     E\u266d major"},
		{&KeySignature{metaevent, -4, 0, "A\u266d major"}, "KeySignature     A\u266d major"},
		{&KeySignature{metaevent, -5, 0, "D\u266d major"}, "KeySignature     D\u266d major"},
		{&KeySignature{metaevent, -6, 0, "G\u266d major"}, "KeySignature     G\u266d major"},

		{&KeySignature{metaevent, 0, 1, "A minor"}, "KeySignature     A minor"},
		{&KeySignature{metaevent, 1, 1, "E minor"}, "KeySignature     E minor"},
		{&KeySignature{metaevent, 2, 1, "B minor"}, "KeySignature     B minor"},
		{&KeySignature{metaevent, 3, 1, "F\u266f minor"}, "KeySignature     F\u266f minor"},
		{&KeySignature{metaevent, 4, 1, "C\u266f minor"}, "KeySignature     C\u266f minor"},
		{&KeySignature{metaevent, 5, 1, "G\u266f minor"}, "KeySignature     G\u266f minor"},
		{&KeySignature{metaevent, 6, 1, "D\u266f minor"}, "KeySignature     D\u266f minor"},
		{&KeySignature{metaevent, -1, 1, "F minor"}, "KeySignature     D minor"},
		{&KeySignature{metaevent, -2, 1, "G minor"}, "KeySignature     G minor"},
		{&KeySignature{metaevent, -3, 1, "C minor"}, "KeySignature     C minor"},
		{&KeySignature{metaevent, -4, 1, "F minor"}, "KeySignature     F minor"},
		{&KeySignature{metaevent, -5, 1, "B\u266d minor"}, "KeySignature     B\u266d minor"},
		{&KeySignature{metaevent, -6, 1, "E\u266d minor"}, "KeySignature     E\u266d minor"},

		{&KeySignature{metaevent, 0, 2, "???"}, "KeySignature     ???"},
	}

	for _, v := range eventlist {
		w := new(bytes.Buffer)

		v.event.Render(w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("KeySignature rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
