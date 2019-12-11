package sysex

import (
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
	"strings"
)

type SysExSingleMessage struct {
	SysExEvent
	Data []byte
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

func (e *SysExSingleMessage) Render(ctx *context.Context, w io.Writer) {
	data := new(bytes.Buffer)

	for _, b := range e.Data {
		fmt.Fprintf(data, "%02X ", b)
	}

	fmt.Fprintf(w, "%s %-16s %s", e.SysExEvent, "SingleMessage", strings.TrimSpace(data.String()))
}
