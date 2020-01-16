package metaevent

import (
	"fmt"
	"io"
)

type CuePoint struct {
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
		MetaEvent: *event,
		CuePoint:  string(cuepoint),
	}, nil
}

func (e *CuePoint) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "CuePoint", e.CuePoint)
}
