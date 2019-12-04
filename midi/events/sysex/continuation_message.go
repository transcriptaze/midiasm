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

func NewSysExContinuationMessage(event *SysExEvent, r io.ByteReader) (*SysExContinuationMessage, error) {
	if event.Status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExContinuationMessage event type (%02x): expected 'F7'", event.Status)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &SysExContinuationMessage{
		SysExEvent: *event,
		Data:       data,
	}, nil
}

func (e *SysExContinuationMessage) Render(ctx *context.Context, w io.Writer) {
	data := new(bytes.Buffer)

	for _, b := range e.Data {
		fmt.Fprintf(data, "%02X ", b)
	}

	fmt.Fprintf(w, "%s %-16s %s", e.SysExEvent, "ContinuationMessage", strings.TrimSpace(data.String()))
}
