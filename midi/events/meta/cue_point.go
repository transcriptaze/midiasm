package metaevent

import (
	"fmt"
	"io"
)

type CuePoint struct {
	Tag string
	MetaEvent
	CuePoint string
}

func NewCuePoint(event *MetaEvent, r io.ByteReader) (*CuePoint, error) {
	if event.Type != 0x07 {
		return nil, fmt.Errorf("Invalid CuePoint event type (%02x): expected '07'", event.Type)
	}

	cuepoint, err := read(r)
	if err != nil {
		return nil, err
	}

	return &CuePoint{
		Tag:       "CuePoint",
		MetaEvent: *event,
		CuePoint:  string(cuepoint),
	}, nil
}
