package midievent

import (
	"fmt"
	"io"
)

type Controller struct {
	MidiEvent
	controller byte
	value      byte
}

func (e *Controller) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d controller:%d value:%d\n", e.Status, "Controller", e.MidiEvent.Event, e.Channel, e.controller, e.value)
}
