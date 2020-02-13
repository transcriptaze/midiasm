package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type TrackName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewTrackName(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*TrackName, error) {
	if eventType != 0x03 {
		return nil, fmt.Errorf("Invalid TrackName event type (%02x): expected '03'", eventType)
	}

	name, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &TrackName{
		Tag:    "TrackName",
		Status: status,
		Type:   eventType,
		Name:   string(name),
	}, nil
}
