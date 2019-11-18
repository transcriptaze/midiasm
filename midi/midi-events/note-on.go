package midievent

import (
	"fmt"
	"io"
)

type NoteOn struct {
	MidiEvent
	note     uint8
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

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d channel:%d note:%d velocity:%d\n", e.Status, "NoteOn", e.Delta, e.channel, e.note, e.velocity)
}
