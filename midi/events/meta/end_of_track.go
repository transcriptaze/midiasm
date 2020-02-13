package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type EndOfTrack struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
}

func NewEndOfTrack(r events.EventReader, status types.Status, eventType types.MetaEventType) (*EndOfTrack, error) {
	if eventType != 0x2f {
		return nil, fmt.Errorf("Invalid EndOfTrack event type (%02x): expected '2f'", eventType)
	}

	data, err := r.ReadVLF()
	if err != nil {
		return nil, err
	} else if len(data) != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", len(data))
	}

	return &EndOfTrack{
		Tag:    "EndOfTrack",
		Status: status,
		Type:   eventType,
	}, nil
}
