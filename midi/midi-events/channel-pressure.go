package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type ChannelPressure struct {
	MidiEvent
	Pressure byte
}

func NewChannelPressure(event MidiEvent, r *bufio.Reader) (*ChannelPressure, error) {
	if event.Status&0xF0 != 0xD0 {
		return nil, fmt.Errorf("Invalid ChannelPressure status (%02x): expected 'D0'", event.Status&0x80)
	}

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	event.bytes = append(event.bytes, pressure)

	return &ChannelPressure{
		MidiEvent: event,
		Pressure:  pressure,
	}, nil
}

func (e *ChannelPressure) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d pressure:%d\n", e.MidiEvent, "ChannelPressure", e.MidiEvent.Event, e.Channel, e.Pressure)
}
