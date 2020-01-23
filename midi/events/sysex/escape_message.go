package sysex

import (
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
	"strings"
)

type SysExEscapeMessage struct {
	SysExEvent
	Data types.Hex
}

func NewSysExEscapeMessage(event *SysExEvent, r io.ByteReader, ctx *context.Context) (*SysExEscapeMessage, error) {
	if event.Status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExEscapeMessage event type (%02x): expected 'F7'", event.Status)
	}

	if ctx.Casio {
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

func (e *SysExEscapeMessage) Render(w io.Writer) {
	data := new(bytes.Buffer)

	for _, b := range e.Data {
		fmt.Fprintf(data, "%02X ", b)
	}

	fmt.Fprintf(w, "%s %-16s %s", e.SysExEvent, "EscapeMessage", strings.TrimSpace(data.String()))
}
