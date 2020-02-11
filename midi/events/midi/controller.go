package midievent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Controller struct {
	Tag string
	*events.Event
	Channel    types.Channel
	Controller byte
	Value      byte
}

func NewController(event *events.Event, r io.ByteReader) (*Controller, error) {
	if event.Status&0xF0 != 0xB0 {
		return nil, fmt.Errorf("Invalid Controller status (%02x): expected 'B0'", event.Status&0x80)
	}

	channel := types.Channel((event.Status) & 0x0F)

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
		Event:      event,
		Channel:    channel,
		Controller: controller,
		Value:      value,
	}, nil
}
