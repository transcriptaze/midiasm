package sysex

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/types"
)

type SysExContinuationMessage struct {
	Tag    string
	Status types.Status
	Data   types.Hex
}

func NewSysExContinuationMessage(ctx *context.Context, r io.ByteReader, status types.Status) (*SysExContinuationMessage, error) {
	if status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExContinuationMessage event type (%02x): expected 'F7'", status)
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		terminator := data[len(data)-1]
		if terminator == 0xf7 {
			data = data[:len(data)-1]
			ctx.Casio = false
		}
	}

	return &SysExContinuationMessage{
		Tag:    "SysExContinuation",
		Status: status,
		Data:   data,
	}, nil
}
