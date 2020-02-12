package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SysExContinuationMessage struct {
	Tag    string
	Status types.Status
	Data   types.Hex
}

func NewSysExContinuationMessage(r io.ByteReader, status types.Status, ctx *context.Context) (*SysExContinuationMessage, error) {
	if status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExContinuationMessage event type (%02x): expected 'F7'", status)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		terminator := data[len(data)-1]
		if terminator == 0xf7 {
			data = data[:len(data)-1]
			ctx.CasioOff()
		}
	}

	return &SysExContinuationMessage{
		Tag:    "SysExContinuation",
		Status: status,
		Data:   data,
	}, nil
}
