package midievent

import (
	"fmt"
	"io"
)

type NoteOff struct {
	delta    uint32
	status   byte
	channel  uint8
	note     uint8
	velocity uint8
	bytes    []byte
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

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d channel:%d note:%d velocity:%d\n", e.status, "NoteOff", e.delta, e.channel, e.note, e.velocity)
}
