package midievent

import (
	"fmt"
	"io"
)

type PolyphonicPressure struct {
	Tag string
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
		Tag:       "PolyphonicPressure",
		MidiEvent: *event,
		Pressure:  pressure,
	}, nil
}
