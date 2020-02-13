package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type Marker struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Marker string
}

func NewMarker(r events.EventReader, status types.Status, eventType types.MetaEventType) (*Marker, error) {
	if eventType != 0x06 {
		return nil, fmt.Errorf("Invalid Marker event type (%02x): expected '06'", eventType)
	}

	marker, err := r.ReadVLF()
	if err != nil {
		return nil, err
	}

	return &Marker{
		Tag:    "Marker",
		Status: status,
		Type:   eventType,
		Marker: string(marker),
	}, nil
}
