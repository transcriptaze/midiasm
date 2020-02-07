package midievent

import (
	"fmt"
	"io"
)

type Controller struct {
	Tag string
	MidiEvent
	Controller byte
	Value      byte
}

func NewController(event *MidiEvent, r io.ByteReader) (*Controller, error) {
	if event.Status&0xF0 != 0xB0 {
		return nil, fmt.Errorf("Invalid Controller status (%02x): expected 'B0'", event.Status&0x80)
	}

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
		MidiEvent:  *event,
		Controller: controller,
		Value:      value,
	}, nil
}
