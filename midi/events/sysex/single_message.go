package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SysExSingleMessage struct {
	SysExEvent
	Manufacturer types.Manufacturer
	Data         types.Hex
}

func NewSysExSingleMessage(event *SysExEvent, r io.ByteReader, ctx *context.Context) (*SysExSingleMessage, error) {
	if event.Status != 0xf0 {
		return nil, fmt.Errorf("Invalid SysExSingleMessage event type (%02x): expected 'F0'", event.Status)
	}

	bytes, err := read(r)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, fmt.Errorf("Invalid SysExSingleMessage event data (0 bytes")
	}

	id := bytes[0:1]
	data := bytes[1:]

	if bytes[0] == 0x00 {
		id = bytes[0:3]
		data = bytes[3:]
	}

	terminator := data[len(data)-1]
	if terminator == 0xf7 {
		data = data[:len(data)-1]
	} else {
		ctx.Casio = true
	}

	return &SysExSingleMessage{
		SysExEvent:   *event,
		Manufacturer: types.LookupManufacturer(id),
		Data:         data,
	}, nil
}
