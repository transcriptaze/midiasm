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
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d pressure:%d\n", e.Status, "ChannelPressure", e.MidiEvent.Event, e.Channel, e.Pressure)
}
