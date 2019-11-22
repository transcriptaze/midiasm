package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type PitchBend struct {
	MidiEvent
	Bend uint16
}

func NewPitchBend(event MidiEvent, r *bufio.Reader) (*PitchBend, error) {
	if event.Status&0xF0 != 0xE0 {
		return nil, fmt.Errorf("Invalid PitchBend status (%02x): expected 'E0'", event.Status&0x80)
	}

	lsb, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	msb, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	event.bytes = append(event.bytes, lsb, msb)

	bend := ((uint16(msb) & 0x7f) << 7) + (uint16(lsb) & 0x7f)

	return &PitchBend{
		MidiEvent: event,
		Bend:      bend,
	}, nil
}

func (e *PitchBend) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d bend:%d\n", e.MidiEvent, "PitchBend", e.Channel, e.Bend)
}
