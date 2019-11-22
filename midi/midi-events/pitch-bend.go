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
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d bend:%d\n", e.Status, "PitchBend", e.MidiEvent.Event, e.Channel, e.Bend)
}
