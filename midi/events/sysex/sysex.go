package sysex

import (
	"fmt"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
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

func Parse(tick uint64, delta uint32, r IO.Reader, status types.Status, ctx *context.Context) (any, error) {
	if status != 0xF0 && status != 0xF7 {
		return nil, fmt.Errorf("Invalid SysEx status (%v): expected 'F0' or 'F7'", status)
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	switch {
	case status == 0xf0 && ctx.Casio:
		return nil, fmt.Errorf("Invalid SysExMessage event data: F0 start byte without terminating F7")

	case status == 0xf0:
		return UnmarshalSysExMessage(ctx, tick, delta, status, data, r.Bytes()...)

	case status == 0xf7 && ctx.Casio:
		return UnmarshalSysExContinuationMessage(ctx, tick, delta, status, data, r.Bytes()...)

	case status == 0xf7:
		return UnmarshalSysExEscapeMessage(ctx, tick, delta, status, data, r.Bytes()...)

	default:
		return nil, fmt.Errorf("Unrecognised SYSEX event: %v", status)
	}
}
