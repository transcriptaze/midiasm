package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SysExEscapeMessage struct {
	SysExEvent
	Data types.Hex
}

func NewSysExEscapeMessage(event *SysExEvent, r io.ByteReader, ctx *context.Context) (*SysExEscapeMessage, error) {
	if event.Status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExEscapeMessage event type (%02x): expected 'F7'", event.Status)
	}

	if ctx.Casio() {
		return nil, fmt.Errorf("F7 is not valid for SysExEscapeMessage event in Casio mode")
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &SysExEscapeMessage{
		SysExEvent: *event,
		Data:       data,
	}, nil
}
