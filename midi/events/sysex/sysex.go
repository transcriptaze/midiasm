package sysex

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TSysExEvent interface {
	SysExMessage | SysExContinuationMessage | SysExEscapeMessage
}

type ISysExEvent interface {
	unmarshal(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error
}

type event struct {
	tick  uint64
	delta lib.Delta
	bytes []byte

	tag    lib.Tag
	Status lib.Status
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
	status := byte(e.Status)

	if delta, err := e.delta.MarshalBinary(); err != nil {
		return nil, err
	} else {
		return append(delta, status), nil
	}
}

func Parse(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	if status != 0xF0 && status != 0xF7 {
		return nil, fmt.Errorf("Invalid SysEx status (%v): expected 'F0' or 'F7'", status)
	}

	switch {
	case status == 0xf0:
		return unmarshal[SysExMessage](tick, delta, status, data, bytes...)

	case status == 0xf7 && ctx.Casio:
		return unmarshal[SysExContinuationMessage](tick, delta, status, data, bytes...)

	case status == 0xf7:
		return unmarshal[SysExEscapeMessage](tick, delta, status, data, bytes...)

	default:
		return nil, fmt.Errorf("Unrecognised SYSEX event: %v", status)
	}
}

// Ref. https://stackoverflow.com/questions/71444847/go-with-generics-type-t-is-pointer-to-type-parameter-not-type-parameter
// Ref. https://stackoverflow.com/questions/69573113/how-can-i-instantiate-a-non-nil-pointer-of-type-argument-with-generic-go/69575720#69575720
func unmarshal[
	E TSysExEvent,
	P interface {
		*E
		ISysExEvent
	}](tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (any, error) {
	p := P(new(E))
	if err := p.unmarshal(tick, delta, status, data, bytes...); err != nil {
		return nil, err
	} else {
		return *p, nil
	}
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
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

func delta(bytes []byte) (lib.Delta, []byte, error) {
	v, remaining, err := vlq(bytes)

	return lib.Delta(v), remaining, err
}
