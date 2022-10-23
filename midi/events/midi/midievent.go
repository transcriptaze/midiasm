package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type event struct {
	tick  uint64
	delta uint32
	bytes []byte

	tag     types.Tag
	Status  types.Status
	Channel types.Channel
}

func (e event) Tick() uint64 {
	return e.tick
}

func (e event) Delta() uint32 {
	return e.delta
}

func (e event) Bytes() []byte {
	return e.bytes
}

func (e event) Tag() string {
	return fmt.Sprintf("%v", e.tag)
}

type Note struct {
	Value byte
	Name  string
	Alias string
}

func Parse(tick uint64, delta uint32, r IO.Reader, status types.Status, ctx *context.Context) (interface{}, error) {
	switch status & 0xF0 {
	case 0x80:
		return UnmarshalNoteOff(ctx, tick, delta, r, status)

	case 0x90:
		return NewNoteOn(ctx, tick, delta, r, status)

	case 0xA0:
		return NewPolyphonicPressure(r, status)

	case 0xB0:
		return NewController(ctx, tick, delta, r, status)

	case 0xC0:
		return NewProgramChange(ctx, tick, delta, r, status)

	case 0xD0:
		return NewChannelPressure(r, status)

	case 0xE0:
		return NewPitchBend(r, status)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
}
