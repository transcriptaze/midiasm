package midievent

import (
	"fmt"
	"io"
)

type ChannelPressure struct {
	Tag string
	MidiEvent
	Pressure byte
}

func NewChannelPressure(event *MidiEvent, r io.ByteReader) (*ChannelPressure, error) {
	if event.Status&0xF0 != 0xD0 {
		return nil, fmt.Errorf("Invalid ChannelPressure status (%02x): expected 'D0'", event.Status&0x80)
	}

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ChannelPressure{
		Tag:       "ChannelPressure",
		MidiEvent: *event,
		Pressure:  pressure,
	}, nil
}
