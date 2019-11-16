package metaevent

import (
	"fmt"
	"io"
)

type TrackName struct {
	MetaEvent
	name string
}

func NewTrackName(event MetaEvent, data []byte) (*TrackName, error) {
	if event.status != 0xff {
		return nil, fmt.Errorf("Invalid TrackName status (%02x): expected 'ff'", event.status)
	}

	if event.eventType != 0x03 {
		return nil, fmt.Errorf("Invalid TrackName event type (%02x): expected '03'", event.eventType)
	}

	name := string(data)

	return &TrackName{
		MetaEvent: event,
		name:      name,
	}, nil
}

func (e *TrackName) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "             ")

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d name:%s\n", e.eventType, "TrackName", e.delta, e.name)
}
