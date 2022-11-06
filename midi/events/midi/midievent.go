package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type event struct {
	tick  uint64
	delta types.Delta
	bytes []byte

	tag     types.Tag
	Status  types.Status
	Channel types.Channel
}

func (e event) Tick() uint64 {
	return e.tick
}

func (e event) Delta() uint32 {
	return uint32(e.delta)
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

var factory = map[byte]func(*context.Context, uint64, uint32, IO.Reader, types.Status) (any, error){
	0x80: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (any, error) {
		return UnmarshalNoteOff(ctx, tick, delta, r, status)
	},

	0x90: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (any, error) {
		return UnmarshalNoteOn(ctx, tick, delta, r, status)
	},

	0xA0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (any, error) {
		return UnmarshalPolyphonicPressure(ctx, tick, delta, r, status)
	},

	0xD0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (any, error) {
		return UnmarshalChannelPressure(ctx, tick, delta, r, status)
	},

	0xE0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (any, error) {
		return UnmarshalPitchBend(ctx, tick, delta, r, status)
	},
}

func Parse(tick uint64, delta uint32, r IO.Reader, status types.Status, ctx *context.Context) (interface{}, error) {
	eventType := byte(status & 0xf0)

	if f, ok := factory[eventType]; ok {
		return f(ctx, tick, delta, r, status)
	}

	switch status & 0xF0 {
	case 0xB0:
		return NewController(ctx, tick, delta, r, status)

	case 0xC0:
		return NewProgramChange(ctx, tick, delta, r, status)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
}
