package midievent

import (
	"fmt"
	"io"
)

type Controller struct {
	delta      uint32
	status     byte
	channel    uint8
	controller byte
	value      byte
	bytes      []byte
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

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d channel:%d controller:%d value:%d\n", e.status, "Controller", e.delta, e.channel, e.controller, e.value)
}
