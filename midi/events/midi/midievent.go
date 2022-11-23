package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

type TMidiEvent interface {
	NoteOff | NoteOn | PolyphonicPressure | Controller | ProgramChange | ChannelPressure | PitchBend
}

type event struct {
	tick  uint64
	delta lib.Delta
	bytes []byte

	tag     lib.Tag
	Status  lib.Status
	Channel lib.Channel
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

func (e event) MarshalBinary() ([]byte, error) {
	status := byte(e.Status & 0xf0)
	channel := byte(e.Channel & 0x0f)

	if delta, err := e.delta.MarshalBinary(); err != nil {
		return nil, err
	} else {
		return append(delta, status|channel), nil
	}
}

type Note struct {
	Value byte
	Name  string
	Alias string
}

var factory = map[byte]func(*context.Context, uint64, uint32, IO.Reader, lib.Status) (any, error){
	0x80: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if note, err := r.ReadByte(); err != nil {
			return nil, err
		} else if velocity, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			return UnmarshalNoteOff(ctx, tick, delta, status, []byte{note, velocity}, r.Bytes()...)
		}
	},

	0x90: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if note, err := r.ReadByte(); err != nil {
			return nil, err
		} else if velocity, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			return UnmarshalNoteOn(ctx, tick, delta, status, []byte{note, velocity}, r.Bytes()...)
		}
	},

	0xA0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if pressure, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			return UnmarshalPolyphonicPressure(tick, delta, status, []byte{pressure}, r.Bytes()...)
		}
	},

	0xB0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if controller, err := r.ReadByte(); err != nil {
			return nil, err
		} else if value, err := r.ReadByte(); err != nil {
			return nil, err
		} else if evt, err := UnmarshalController(tick, delta, status, []byte{controller, value}, r.Bytes()...); err != nil {
			return nil, err
		} else {
			if ctx != nil && evt.Controller.ID == 0x00 {
				c := uint8(evt.Channel)
				v := uint16(evt.Value)
				ctx.ProgramBank[c] = (ctx.ProgramBank[c] & 0x003f) | ((v & 0x003f) << 7)
			}

			if ctx != nil && evt.Controller.ID == 0x20 {
				c := uint8(evt.Channel)
				v := uint16(evt.Value)
				ctx.ProgramBank[c] = (ctx.ProgramBank[c] & (0x003f << 7)) | (v & 0x003f)
			}

			return evt, err
		}
	},

	0xC0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if program, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			return UnmarshalProgramChange(ctx, tick, delta, status, []byte{program}, r.Bytes()...)
		}
	},

	0xD0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if pressure, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			return UnmarshalChannelPressure(tick, delta, status, []byte{pressure}, r.Bytes()...)
		}
	},

	0xE0: func(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (any, error) {
		if msb, err := r.ReadByte(); err != nil {
			return nil, err
		} else if lsb, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			return UnmarshalPitchBend(tick, delta, status, []byte{msb, lsb}, r.Bytes()...)
		}
	},
}

func Parse(tick uint64, delta uint32, r IO.Reader, status lib.Status, ctx *context.Context) (interface{}, error) {
	eventType := byte(status & 0xf0)

	if f, ok := factory[eventType]; ok {
		return f(ctx, tick, delta, r, status)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
}

func or(status lib.Status, channel lib.Channel) lib.Status {
	return lib.Status(byte(status&0xf0) | byte(channel&0x0f))
}
