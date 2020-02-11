package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type PitchBend struct {
	Tag string
	*events.Event
	Channel types.Channel
	Bend    uint16
}

func NewPitchBend(event *events.Event, r io.ByteReader) (*PitchBend, error) {
	if event.Status&0xF0 != 0xE0 {
		return nil, fmt.Errorf("Invalid PitchBend status (%02x): expected 'E0'", event.Status&0x80)
	}

	channel := types.Channel((event.Status) & 0x0F)

	bend := uint16(0)

	for i := 0; i < 2; i++ {
		if b, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			bend <<= 7
			bend |= uint16(b) & 0x7f
		}
	}

	return &PitchBend{
		Tag:     "PitchBend",
		Event:   event,
		Channel: channel,
		Bend:    bend,
	}, nil
}
