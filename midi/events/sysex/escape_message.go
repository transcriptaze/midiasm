package sysex

import (
	"fmt"
	"io"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/types"
)

type SysExEscapeMessage struct {
	Tag    string
	Status types.Status
	Data   types.Hex
}

func NewSysExEscapeMessage(ctx *context.Context, r io.ByteReader, status types.Status) (*SysExEscapeMessage, error) {
	if status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExEscapeMessage event type (%02x): expected 'F7'", status)
	}

	if ctx.Casio {
		return nil, fmt.Errorf("F7 is not valid for SysExEscapeMessage event in Casio mode")
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &SysExEscapeMessage{
		Tag:    "SysExEscape",
		Status: status,
		Data:   data,
	}, nil
}
