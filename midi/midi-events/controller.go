package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type Controller struct {
	MidiEvent
	Controller byte
	Value      byte
}

func NewController(event MidiEvent, r *bufio.Reader) (*Controller, error) {
	if event.Status&0xF0 != 0xB0 {
		return nil, fmt.Errorf("Invalid Controller status (%02x): expected 'B0'", event.Status&0x80)
	}

	controller, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	value, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	event.bytes = append(event.bytes, controller, value)

	return &Controller{
		MidiEvent:  event,
		Controller: controller,
		Value:      value,
	}, nil
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

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d controller:%d value:%d\n", e.Status, "Controller", e.MidiEvent.Event, e.Channel, e.Controller, e.Value)
}
