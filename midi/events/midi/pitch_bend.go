package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type PitchBend struct {
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
		MidiEvent: *event,
		Bend:      bend,
	}, nil
}

func (e *PitchBend) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d bend:%d", e.MidiEvent, "PitchBend", e.Channel, e.Bend)
}
