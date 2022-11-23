package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
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

var EVENTS = map[byte]int{
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

type Note struct {
	Value byte
	Name  string
	Alias string
}

func Parse(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	return parse(ctx, tick, delta, status, data, bytes...)
}

func parse(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	switch status & 0xf0 {
	case 0x80:
		e, err := UnmarshalNoteOff(ctx, tick, delta, status, data)
		if e != nil && err == nil {
			e.bytes = bytes
		}
		return e, err

	case 0x90:
		return UnmarshalNoteOn(ctx, tick, delta, status, data, bytes...)

	case 0xA0:
		return UnmarshalPolyphonicPressure(tick, delta, status, data, bytes...)

	case 0xB0:
		if evt, err := UnmarshalController(tick, delta, status, data, bytes...); err != nil {
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

	case 0xC0:
		return UnmarshalProgramChange(ctx, tick, delta, status, data, bytes...)

	case 0xD0:
		return UnmarshalChannelPressure(tick, delta, status, data, bytes...)

	case 0xE0:
		return UnmarshalPitchBend(tick, delta, status, data, bytes...)

	default:
		return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
	}
}

func or(status lib.Status, channel lib.Channel) lib.Status {
	return lib.Status(byte(status&0xf0) | byte(channel&0x0f))
}
