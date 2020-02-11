package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type ChannelPressure struct {
	Tag string
	*events.Event
	Channel  types.Channel
	Pressure byte
}

func NewChannelPressure(event *events.Event, r io.ByteReader) (*ChannelPressure, error) {
	if event.Status&0xF0 != 0xD0 {
		return nil, fmt.Errorf("Invalid ChannelPressure status (%02x): expected 'D0'", event.Status&0x80)
	}

	channel := types.Channel((event.Status) & 0x0F)

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ChannelPressure{
		Tag:      "ChannelPressure",
		Event:    event,
		Channel:  channel,
		Pressure: pressure,
	}, nil
}
