package midievent

import (
	"fmt"
	"io"
)

type NoteOn struct {
	MidiEvent
	Note     uint8
	velocity uint8
}

func (e *NoteOn) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d note:%d velocity:%d\n", e.Status, "NoteOn", e.MidiEvent.Event, e.Channel, e.Note, e.velocity)
}
