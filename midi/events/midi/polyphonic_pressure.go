package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type PolyphonicPressure struct {
	Tag      string
	Status   types.Status
	Channel  types.Channel
	Pressure byte
}

func NewPolyphonicPressure(r io.ByteReader, status types.Status) (*PolyphonicPressure, error) {
	if status&0xF0 != 0xA0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%v): expected 'Ax'", status)
	}

	channel := types.Channel(status & 0x0F)

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &PolyphonicPressure{
		Tag:      "PolyphonicPressure",
		Status:   status,
		Channel:  channel,
		Pressure: pressure,
	}, nil
}
