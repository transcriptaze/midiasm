package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"io"
)

type Note struct {
	Value byte
	Name  string
	Alias string
}

func Parse(event *events.Event, r io.ByteReader, ctx *context.Context) (interface{}, error) {
	switch event.Status & 0xF0 {
	case 0x80:
		return NewNoteOff(ctx, event, r)

	case 0x90:
		return NewNoteOn(ctx, event, r)

	case 0xA0:
		return NewPolyphonicPressure(event, r)

	case 0xB0:
		return NewController(event, r)

	case 0xC0:
		return NewProgramChange(event, r)

	case 0xD0:
		return NewChannelPressure(event, r)

	case 0xE0:
		return NewPitchBend(event, r)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %02X", byte(event.Status&0xF0))
}
