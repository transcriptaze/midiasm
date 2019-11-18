package midievent

import (
	"fmt"
	"io"
)

type ProgramChange struct {
	MidiEvent
	program uint8
}

func (e *ProgramChange) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d channel:%d program:%d\n", e.Status, "ProgramChange", e.Delta, e.channel, e.program)
}
