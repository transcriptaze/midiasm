package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type CuePoint struct {
	Tag      string
	Status   types.Status
	Type     types.MetaEventType
	CuePoint string
}

func NewCuePoint(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*CuePoint, error) {
	if eventType != 0x07 {
		return nil, fmt.Errorf("Invalid CuePoint event type (%02x): expected '07'", eventType)
	}

	cuepoint, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &CuePoint{
		Tag:      "CuePoint",
		Status:   status,
		Type:     eventType,
		CuePoint: string(cuepoint),
	}, nil
}
