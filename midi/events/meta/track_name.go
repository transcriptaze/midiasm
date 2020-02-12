package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type TrackName struct {
	Tag string
	MetaEvent
	Type types.MetaEventType
	Name string
}

func NewTrackName(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*TrackName, error) {
	if eventType != 0x03 {
		return nil, fmt.Errorf("Invalid TrackName event type (%02x): expected '03'", eventType)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &TrackName{
		Tag:       "TrackName",
		MetaEvent: *event,
		Type:      eventType,
		Name:      string(name),
	}, nil
}
