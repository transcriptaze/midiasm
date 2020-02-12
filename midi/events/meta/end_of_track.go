package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type EndOfTrack struct {
	Tag string
	MetaEvent
	Type types.MetaEventType
}

func NewEndOfTrack(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*EndOfTrack, error) {
	if eventType != 0x2f {
		return nil, fmt.Errorf("Invalid EndOfTrack event type (%02x): expected '2f'", eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", len(data))
	}

	return &EndOfTrack{
		Tag:       "EndOfTrack",
		Type:      eventType,
		MetaEvent: *event,
	}, nil
}
