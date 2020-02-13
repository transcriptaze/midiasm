package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type DeviceName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewDeviceName(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*DeviceName, error) {
	if eventType != 0x09 {
		return nil, fmt.Errorf("Invalid DeviceName event type (%02x): expected '09'", eventType)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &DeviceName{
		Tag:    "DeviceName",
		Status: status,
		Type:   eventType,
		Name:   string(name),
	}, nil
}
