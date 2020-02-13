package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Marker struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Marker string
}

func NewMarker(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*Marker, error) {
	if eventType != 0x06 {
		return nil, fmt.Errorf("Invalid Marker event type (%02x): expected '06'", eventType)
	}

	marker, err := read(r)
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
