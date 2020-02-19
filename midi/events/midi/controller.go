package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Controller struct {
	Tag        string
	Status     types.Status
	Channel    types.Channel
	Controller types.Controller
	Value      byte
}

func NewController(ctx *context.Context, r io.ByteReader, status types.Status) (*Controller, error) {
	if status&0xF0 != 0xB0 {
		return nil, fmt.Errorf("Invalid Controller status (%v): expected 'Bx'", status)
	}

	channel := uint8(status & 0x0f)

	controller, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	value, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if controller == 0x00 {
		ctx.ProgramBank[channel] = (ctx.ProgramBank[channel] & 0x003f) | ((uint16(value) & 0x003f) << 7)
	}

	if controller == 0x20 {
		ctx.ProgramBank[channel] = (ctx.ProgramBank[channel] & (0x003f << 7)) | (uint16(value) & 0x003f)
	}

	return &Controller{
		Tag:        "Controller",
		Status:     status,
		Channel:    types.Channel(channel),
		Controller: types.LookupController(controller),
		Value:      value,
	}, nil
}
