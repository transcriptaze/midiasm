package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type PolyphonicPressure struct {
	MidiEvent
	Pressure byte
}

func NewPolyphonicPressure(event MidiEvent, r *bufio.Reader) (*PolyphonicPressure, error) {
	if event.Status&0xF0 != 0xA0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%02x): expected 'A0'", event.Status&0x80)
	}

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	event.bytes = append(event.bytes, pressure)

	return &PolyphonicPressure{
		MidiEvent: event,
		Pressure:  pressure,
	}, nil
}

func (e *PolyphonicPressure) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d pressure:%d\n", e.Status, "PolyphonicPressure", e.MidiEvent.Event, e.Channel, e.Pressure)
}
