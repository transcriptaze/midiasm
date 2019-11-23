package midievent

import (
	"fmt"
	"io"
)

type PolyphonicPressure struct {
	MidiEvent
	Pressure byte
}

func NewPolyphonicPressure(event *MidiEvent, r io.ByteReader) (*PolyphonicPressure, error) {
	if event.Status&0xF0 != 0xA0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%02x): expected 'A0'", event.Status&0x80)
	}

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &PolyphonicPressure{
		MidiEvent: *event,
		Pressure:  pressure,
	}, nil
}

func (e *PolyphonicPressure) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d pressure:%d\n", e.MidiEvent, "PolyphonicPressure", e.Channel, e.Pressure)
}
