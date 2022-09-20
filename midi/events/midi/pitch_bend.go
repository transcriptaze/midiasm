package midievent

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/types"
)

type PitchBend struct {
	Tag     string
	Status  types.Status
	Channel types.Channel
	Bend    uint16
}

func NewPitchBend(r io.ByteReader, status types.Status) (*PitchBend, error) {
	if status&0xF0 != 0xE0 {
		return nil, fmt.Errorf("Invalid PitchBend status (%v): expected 'Ex'", status)
	}

	channel := types.Channel(status & 0x0F)

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
		Status:  status,
		Channel: channel,
		Bend:    bend,
	}, nil
}
