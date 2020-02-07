package midievent

import (
	"fmt"
	"io"
)

type PitchBend struct {
	Tag string
	MidiEvent
	Bend uint16
}

func NewPitchBend(event *MidiEvent, r io.ByteReader) (*PitchBend, error) {
	if event.Status&0xF0 != 0xE0 {
		return nil, fmt.Errorf("Invalid PitchBend status (%02x): expected 'E0'", event.Status&0x80)
	}

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
		Tag:       "PitchBend",
		MidiEvent: *event,
		Bend:      bend,
	}, nil
}
