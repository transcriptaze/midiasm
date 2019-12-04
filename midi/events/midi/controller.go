package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type Controller struct {
	MidiEvent
	Controller byte
	Value      byte
}

func NewController(event *MidiEvent, r io.ByteReader) (*Controller, error) {
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

	return &Controller{
		MidiEvent:  *event,
		Controller: controller,
		Value:      value,
	}, nil
}

func (e *Controller) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d controller:%d value:%d", e.MidiEvent, "Controller", e.Channel, e.Controller, e.Value)
}
