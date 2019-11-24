package metaevent

import (
	"fmt"
	"io"
)

type EndOfTrack struct {
	MetaEvent
}

func NewEndOfTrack(event *MetaEvent, r io.ByteReader) (*EndOfTrack, error) {
	if event.eventType != 0x2f {
		return nil, fmt.Errorf("Invalid EndOfTrack event type (%02x): expected '2f'", event.eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", len(data))
	}

	return &EndOfTrack{
		MetaEvent: *event,
	}, nil
}

func (e *EndOfTrack) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s\n", e.MetaEvent, "EndOfTrack")
}
