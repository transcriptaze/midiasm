package sysex

import (
	"bytes"
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
	"strings"
)

type SysExContinuationMessage struct {
	SysExEvent
	Data []byte
}

func NewSysExContinuationMessage(event *SysExEvent, r io.ByteReader, ctx *context.Context) (*SysExContinuationMessage, error) {
	if event.Status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExContinuationMessage event type (%02x): expected 'F7'", event.Status)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		terminator := data[len(data)-1]
		if terminator == 0xf7 {
			ctx.Casio = false
		}
	}

	return &SysExContinuationMessage{
		SysExEvent: *event,
		Data:       data,
	}, nil
}

func (e *SysExContinuationMessage) Render(w io.Writer) {
	data := new(bytes.Buffer)

	for _, b := range e.Data {
		fmt.Fprintf(data, "%02X ", b)
	}

	fmt.Fprintf(w, "%s %-16s %s", e.SysExEvent, "ContinuationMessage", strings.TrimSpace(data.String()))
}
