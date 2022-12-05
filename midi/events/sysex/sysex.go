package sysex

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type TSysExEvent interface {
	SysExMessage | SysExContinuationMessage | SysExEscapeMessage
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
	case status == 0xf0 && ctx.Casio:
		return nil, fmt.Errorf("Invalid SysExMessage event data: F0 start byte without terminating F7")

	case status == 0xf0:
		if e, err := UnmarshalSysExMessage(ctx, tick, delta, status, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case status == 0xf7 && ctx.Casio:
		if e, err := UnmarshalSysExContinuationMessage(ctx, tick, delta, status, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	case status == 0xf7:
		if e, err := UnmarshalSysExEscapeMessage(ctx, tick, delta, status, data, bytes...); err != nil || e == nil {
			return nil, err
		} else {
			return *e, err
		}

	default:
		return nil, fmt.Errorf("Unrecognised SYSEX event: %v", status)
	}
}

func equal(s string, tag lib.Tag) bool {
	return s == fmt.Sprintf("%v", tag)
}
