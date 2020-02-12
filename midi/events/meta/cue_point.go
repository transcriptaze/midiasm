package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type CuePoint struct {
	Tag string
	MetaEvent
	Type     types.MetaEventType
	CuePoint string
}

func NewCuePoint(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*CuePoint, error) {
	if eventType != 0x07 {
		return nil, fmt.Errorf("Invalid CuePoint event type (%02x): expected '07'", eventType)
	}

	cuepoint, err := read(r)
	if err != nil {
		return nil, err
	}

	return &CuePoint{
		Tag:       "CuePoint",
		MetaEvent: *event,
		Type:      eventType,
		CuePoint:  string(cuepoint),
	}, nil
}
