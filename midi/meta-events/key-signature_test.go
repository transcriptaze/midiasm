package metaevent

import (
	"bytes"
	"github.com/twystd/midiasm/midi/event"
	"strings"
	"testing"
)

func TestKeySignatureRender(t *testing.T) {
	ctx := event.Context{
		Scale: event.Sharps,
	}

	metaevent := MetaEvent{
		event.Event{76, 12, 0xff, []byte{0x00, 0xff, 0x59, 0x02, 0x00, 0x00}},
		0x59,
	}

	var events = []struct {
		event    event.IEvent
		expected string
	}{
		{&KeySignature{metaevent, 0, 0}, "KeySignature     C major"},
		{&KeySignature{metaevent, 1, 0}, "KeySignature     G major"},
		{&KeySignature{metaevent, 2, 0}, "KeySignature     D major"},
		{&KeySignature{metaevent, 3, 0}, "KeySignature     A major"},
		{&KeySignature{metaevent, 4, 0}, "KeySignature     E major"},
		{&KeySignature{metaevent, 5, 0}, "KeySignature     B major"},
		{&KeySignature{metaevent, 6, 0}, "KeySignature     F\u266f major"},
		{&KeySignature{metaevent, -1, 0}, "KeySignature     F major"},
		{&KeySignature{metaevent, -2, 0}, "KeySignature     B\u266d major"},
		{&KeySignature{metaevent, -3, 0}, "KeySignature     E\u266d major"},
		{&KeySignature{metaevent, -4, 0}, "KeySignature     A\u266d major"},
		{&KeySignature{metaevent, -5, 0}, "KeySignature     D\u266d major"},
		{&KeySignature{metaevent, -6, 0}, "KeySignature     G\u266d major"},

		{&KeySignature{metaevent, 0, 1}, "KeySignature     A minor"},
		{&KeySignature{metaevent, 1, 1}, "KeySignature     E minor"},
		{&KeySignature{metaevent, 2, 1}, "KeySignature     B minor"},
		{&KeySignature{metaevent, 3, 1}, "KeySignature     F\u266f minor"},
		{&KeySignature{metaevent, 4, 1}, "KeySignature     C\u266f minor"},
		{&KeySignature{metaevent, 5, 1}, "KeySignature     G\u266f minor"},
		{&KeySignature{metaevent, 6, 1}, "KeySignature     D\u266f minor"},
		{&KeySignature{metaevent, -1, 1}, "KeySignature     D minor"},
		{&KeySignature{metaevent, -2, 1}, "KeySignature     G minor"},
		{&KeySignature{metaevent, -3, 1}, "KeySignature     C minor"},
		{&KeySignature{metaevent, -4, 1}, "KeySignature     F minor"},
		{&KeySignature{metaevent, -5, 1}, "KeySignature     B\u266d minor"},
		{&KeySignature{metaevent, -6, 1}, "KeySignature     E\u266d minor"},

		{&KeySignature{metaevent, 0, 2}, "KeySignature     ???"},
	}

	for _, v := range events {
		w := new(bytes.Buffer)

		v.event.Render(&ctx, w)

		if !strings.HasSuffix(w.String(), v.expected) {
			t.Errorf("KeySignature rendered incorrectly\nExpected: '%s'\ngot:      '%s'", v.expected, w.String())
		}
	}
}
