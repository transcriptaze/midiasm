package midievent

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TMidiEvent interface {
	NoteOff | NoteOn | PolyphonicPressure | Controller | ProgramChange | ChannelPressure | PitchBend
}

type IMidiEvent interface {
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
	Value byte   `json:"value"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
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
		return unmarshal[NoteOff](ctx, tick, delta, status, data, bytes...)

	case 0x90:
		return unmarshal[NoteOn](ctx, tick, delta, status, data, bytes...)

	case 0xA0:
		return unmarshal[PolyphonicPressure](ctx, tick, delta, status, data, bytes...)

	case 0xB0:
		return unmarshal[Controller](ctx, tick, delta, status, data, bytes...)

	case 0xC0:
		return unmarshal[ProgramChange](ctx, tick, delta, status, data, bytes...)

	case 0xD0:
		return unmarshal[ChannelPressure](ctx, tick, delta, status, data, bytes...)

	case 0xE0:
		return unmarshal[PitchBend](ctx, tick, delta, status, data, bytes...)
	}

	return nil, fmt.Errorf("Unrecognised MIDI event: %v", status)
}

// Ref. https://stackoverflow.com/questions/71444847/go-with-generics-type-t-is-pointer-to-type-parameter-not-type-parameter
// Ref. https://stackoverflow.com/questions/69573113/how-can-i-instantiate-a-non-nil-pointer-of-type-argument-with-generic-go/69575720#69575720
func unmarshal[
	E TMidiEvent,
	P interface {
		*E
		IMidiEvent
	}](ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	p := P(new(E))
	if err := p.unmarshal(ctx, tick, delta, status, data, bytes...); err != nil {
		return nil, err
	} else {
		return *p, nil
	}
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}

func or(status lib.Status, channel lib.Channel) lib.Status {
	return lib.Status(byte(status&0xf0) | byte(channel&0x0f))
}
