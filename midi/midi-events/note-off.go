package midievent

import (
	"fmt"
	"io"
)

type NoteOff struct {
	MidiEvent
	note     uint8
	velocity uint8
}

func (e *NoteOff) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d note:%d velocity:%d\n", e.Status, "NoteOff", e.MidiEvent.Event, e.channel, e.note, e.velocity)
}
