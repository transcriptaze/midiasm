package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type Copyright struct {
	Tag       string
	Status    types.Status
	Type      types.MetaEventType
	Copyright string
}

func NewCopyright(r events.EventReader, status types.Status, eventType types.MetaEventType) (*Copyright, error) {
	if eventType != 0x02 {
		return nil, fmt.Errorf("Invalid Copyright event type (%02x): expected '02'", eventType)
	}

	data, err := r.ReadVLQ()
	if err != nil {
		return nil, err
	}

	return &Copyright{
		Tag:       "Copyright",
		Status:    status,
		Type:      eventType,
		Copyright: string(data),
	}, nil
}
