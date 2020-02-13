package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type Controller struct {
	Tag        string
	Status     types.Status
	Channel    types.Channel
	Controller byte
	Value      byte
}

func NewController(r events.EventReader, status types.Status) (*Controller, error) {
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
		Controller: controller,
		Value:      value,
	}, nil
}
