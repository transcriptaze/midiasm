package sysex

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

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

func Parse(tick uint64, delta uint32, r IO.Reader, status types.Status, ctx *context.Context) (interface{}, error) {
	if status != 0xF0 && status != 0xF7 {
		return nil, fmt.Errorf("Invalid SysEx tag (%v): expected 'F0' or 'F7'", status)
	}

	switch status {
	case 0xf0:
		if ctx.Casio {
			return nil, fmt.Errorf("Invalid SysExSingleMessage event data: F0 start byte without terminating F7")
		} else {
			return UnmarshalSysExSingleMessage(ctx, tick, delta, r, status)
		}

	case 0xf7:
		if ctx.Casio {
			return NewSysExContinuationMessage(ctx, r, status)
		} else {
			return NewSysExEscapeMessage(ctx, r, status)
		}
	}

	return nil, fmt.Errorf("Unrecognised SYSEX event: %v", status)
}
