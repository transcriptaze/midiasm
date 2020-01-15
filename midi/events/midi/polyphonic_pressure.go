package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type PolyphonicPressure struct {
	MidiEvent
	Pressure byte
}

func NewPolyphonicPressure(event *MidiEvent, r io.ByteReader) (*PolyphonicPressure, error) {
	if event.Status&0xF0 != 0xA0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%02x): expected 'A0'", event.Status&0x80)
	}

	pressure, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &PolyphonicPressure{
		MidiEvent: *event,
		Pressure:  pressure,
	}, nil
}

func (e *PolyphonicPressure) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%-2v pressure:%d", e.MidiEvent, "PolyphonicPressure", e.Channel, e.Pressure)
}
