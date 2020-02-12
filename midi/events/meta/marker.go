package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Marker struct {
	Tag string
	*events.Event
	Type   types.MetaEventType
	Marker string
}

func NewMarker(event *events.Event, eventType types.MetaEventType, r io.ByteReader) (*Marker, error) {
	if eventType != 0x06 {
		return nil, fmt.Errorf("Invalid Marker event type (%02x): expected '06'", eventType)
	}

	marker, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Marker{
		Tag:    "Marker",
		Event:  event,
		Type:   eventType,
		Marker: string(marker),
	}, nil
}
