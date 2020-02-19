package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

func Parse(r io.ByteReader, status types.Status, ctx *context.Context) (interface{}, error) {
	if status != 0xF0 && status != 0xF7 {
		return nil, fmt.Errorf("Invalid SysEx tag (%v): expected 'F0' or 'F7'", status)
	}

	switch status {
	case 0xf0:
		if ctx.Casio {
			return nil, fmt.Errorf("Invalid SysExSingleMessage event data: F0 start byte without terminating F7")
		} else {
			return NewSysExSingleMessage(ctx, r, status)
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
