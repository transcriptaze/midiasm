package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

func Parse(r io.ByteReader, status types.Status, ctx *context.Context) (interface{}, error) {
	switch status & 0xF0 {
	case 0x80:
		return NewNoteOff(ctx, r, status)

	case 0x90:
		return NewNoteOn(ctx, r, status)

	case 0xA0:
		return NewPolyphonicPressure(r, status)

	case 0xB0:
		return NewController(r, status)

	case 0xC0:
		return NewProgramChange(r, status)

	case 0xD0:
		return NewChannelPressure(r, status)

	case 0xE0:
		return NewPitchBend(r, status)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
}
