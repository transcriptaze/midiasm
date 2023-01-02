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
	unmarshal(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error
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

func Parse(tick uint64, runningStatus byte, bytes ...byte) (any, error) {
	var status lib.Status
	var data []byte

	delta, remaining, err := vlq(bytes)
	if err != nil {
		return nil, err
	} else if len(remaining) < 1 {
		return nil, fmt.Errorf("Invalid MIDI event - missing status")
	}

	if remaining[0] < 0x80 {
		status = lib.Status(runningStatus)
		data = remaining
	} else {
		status = lib.Status(remaining[0])
		data = remaining[1:]
	}

	switch status & 0xf0 {
	case 0x80:
		return unmarshal[NoteOff](tick, delta, status, data, bytes...)

	case 0x90:
		return unmarshal[NoteOn](tick, delta, status, data, bytes...)

	case 0xA0:
		return unmarshal[PolyphonicPressure](tick, delta, status, data, bytes...)

	case 0xB0:
		return unmarshal[Controller](tick, delta, status, data, bytes...)

	case 0xC0:
		return unmarshal[ProgramChange](tick, delta, status, data, bytes...)

	case 0xD0:
		return unmarshal[ChannelPressure](tick, delta, status, data, bytes...)

	case 0xE0:
		return unmarshal[PitchBend](tick, delta, status, data, bytes...)
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
	}](tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	p := P(new(E))
	if err := p.unmarshal(tick, delta, status, data, bytes...); err != nil {
		return nil, err
	} else {
		return *p, nil
	}
}

func vlq(bytes []byte) (uint32, []byte, error) {
	vlq := uint32(0)

	for i, b := range bytes {
		vlq <<= 7
		vlq += uint32(b & 0x7f)

		if b&0x80 == 0 {
			return vlq, bytes[i+1:], nil
		}
	}

	return 0, nil, fmt.Errorf("Invalid event 'delta'")
}

func vlf(bytes []byte) ([]byte, error) {
	if N, remaining, err := vlq(bytes); err != nil {
		return nil, err
	} else {
		return remaining[:N], nil
	}
}

func delta(bytes []byte) (uint32, []byte, error) {
	v, remaining, err := vlq(bytes)

	return v, remaining, err
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}

func or(status lib.Status, channel lib.Channel) lib.Status {
	return lib.Status(byte(status&0xf0) | byte(channel&0x0f))
}

func FormatNote(ctx *context.Context, n byte) string {
	if ctx != nil {
		return ctx.FormatNote(n)
	}

	var scale = context.Sharps
	var note = scale[n%12]
	var octave int

	if context.MiddleC == lib.C4 {
		octave = int(n/12) - 2
	} else {
		octave = int(n/12) - 1
	}

	return fmt.Sprintf("%s%d", note, octave)
}
