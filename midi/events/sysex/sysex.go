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
		if ctx.Casio() {
			return nil, fmt.Errorf("Invalid SysExSingleMessage event data: F0 start byte without terminating F7")
		} else {
			return NewSysExSingleMessage(ctx, status, r)
		}

	case 0xf7:
		if ctx.Casio() {
			return NewSysExContinuationMessage(r, status, ctx)
		} else {
			return NewSysExEscapeMessage(r, status, ctx)
		}
	}

	return nil, fmt.Errorf("Unrecognised SYSEX event: %v", status)
}

func read(r io.ByteReader) ([]byte, error) {
	N, err := vlq(r)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, N)

	for i := 0; i < int(N); i++ {
		if b, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			bytes[i] = b
		}
	}

	return bytes, nil
}

func vlq(r io.ByteReader) (uint32, error) {
	l := uint32(0)

	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		l <<= 7
		l += uint32(b & 0x7f)

		if b&0x80 == 0 {
			break
		}
	}

	return l, nil
}
