package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SysExSingleMessage struct {
	SysExEvent
	Data types.Hex
}

func NewSysExSingleMessage(event *SysExEvent, r io.ByteReader, ctx *context.Context) (*SysExSingleMessage, error) {
	if event.Status != 0xf0 {
		return nil, fmt.Errorf("Invalid SysExSingleMessage event type (%02x): expected 'F0'", event.Status)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("Invalid SysExSingleMessage event data (0 bytes")
	}

	terminator := data[len(data)-1]
	if terminator != 0xf7 {
		ctx.Casio = true
	}

	return &SysExSingleMessage{
		SysExEvent: *event,
		Data:       data,
	}, nil
}
