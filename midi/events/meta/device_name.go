package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type DeviceName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewDeviceName(r events.EventReader, status types.Status, eventType types.MetaEventType) (*DeviceName, error) {
	if eventType != 0x09 {
		return nil, fmt.Errorf("Invalid DeviceName event type (%02x): expected '09'", eventType)
	}

	name, err := r.ReadVLQ()
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
