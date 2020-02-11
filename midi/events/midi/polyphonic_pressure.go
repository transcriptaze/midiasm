package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type PolyphonicPressure struct {
	Tag string
	*events.Event
	Channel  types.Channel
	Pressure byte
}

func NewPolyphonicPressure(event *events.Event, r io.ByteReader) (*PolyphonicPressure, error) {
	if event.Status&0xF0 != 0xA0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%02x): expected 'A0'", event.Status&0x80)
	}

	channel := types.Channel((event.Status) & 0x0F)

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &PolyphonicPressure{
		Tag:      "PolyphonicPressure",
		Event:    event,
		Channel:  channel,
		Pressure: pressure,
	}, nil
}
