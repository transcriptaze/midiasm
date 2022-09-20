package midievent

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/types"
)

type ChannelPressure struct {
	Tag      string
	Status   types.Status
	Channel  types.Channel
	Pressure byte
}

func NewChannelPressure(r io.ByteReader, status types.Status) (*ChannelPressure, error) {
	if status&0xF0 != 0xD0 {
		return nil, fmt.Errorf("Invalid ChannelPressure status (%v): expected 'Dx'", status)
	}

	channel := types.Channel(status & 0x0F)

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ChannelPressure{
		Tag:      "ChannelPressure",
		Status:   status,
		Channel:  channel,
		Pressure: pressure,
	}, nil
}
