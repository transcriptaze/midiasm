package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TMidiEvent interface {
	NoteOff | NoteOn | PolyphonicPressure | Controller | ProgramChange | ChannelPressure | PitchBend
}

type TMidiEventQ interface {
	*NoteOff | *NoteOn | *PolyphonicPressure | *Controller | *ProgramChange | *ChannelPressure | *PitchBend

	unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error
}

type event struct {
	tick  uint64
	delta lib.Delta
	bytes []byte

	tag     lib.Tag
	Status  lib.Status
	Channel lib.Channel
}

type Note struct {
	Value byte
	Name  string
	Alias string
}

var Events = map[byte]int{
	0x80: 2,
	0x90: 2,
	0xA0: 1,
	0xB0: 2,
	0xC0: 1,
	0xD0: 1,
	0xE0: 2,
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

func Parse(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	switch status & 0xf0 {
	case 0x80:
		e := NoteOff{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}

	case 0x90:
		e := NoteOn{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}

	case 0xA0:
		e := PolyphonicPressure{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}

	case 0xB0:
		e := Controller{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}

	case 0xC0:
		e := ProgramChange{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}

	case 0xD0:
		e := ChannelPressure{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}

	case 0xE0:
		e := PitchBend{}
		if err := parse(&e, ctx, tick, delta, status, data, bytes...); err != nil {
			return nil, err
		} else {
			return e, err
		}
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
}

func parse[T TMidiEventQ](e T, ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	return e.unmarshal(ctx, tick, delta, status, data, bytes...)
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}

func or(status lib.Status, channel lib.Channel) lib.Status {
	return lib.Status(byte(status&0xf0) | byte(channel&0x0f))
}
