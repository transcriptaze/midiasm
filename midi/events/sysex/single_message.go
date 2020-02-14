package sysex

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SysExSingleMessage struct {
	Tag          string
	Status       types.Status
	Manufacturer types.Manufacturer
	Data         types.Hex
}

func NewSysExSingleMessage(ctx *context.Context, r io.ByteReader, status types.Status) (*SysExSingleMessage, error) {
	if status != 0xf0 {
		return nil, fmt.Errorf("Invalid SysExSingleMessage event type (%02x): expected 'F0'", status)
	}

	bytes, err := events.VLF(r)
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
		ctx.CasioOn()
	}

	return &SysExSingleMessage{
		Tag:          "SysExMessage",
		Status:       status,
		Manufacturer: types.LookupManufacturer(id),
		Data:         data,
	}, nil
}
