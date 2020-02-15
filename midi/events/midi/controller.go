package midievent

import (
	"fmt"
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

func NewController(r io.ByteReader, status types.Status) (*Controller, error) {
	if status&0xF0 != 0xB0 {
		return nil, fmt.Errorf("Invalid Controller status (%v): expected 'Bx'", status)
	}

	channel := types.Channel(status & 0x0F)

	controller, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	value, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &Controller{
		Tag:        "Controller",
		Status:     status,
		Channel:    channel,
		Controller: types.LookupController(controller),
		Value:      value,
	}, nil
}
