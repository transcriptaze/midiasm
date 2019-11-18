package metaevent

import (
	"fmt"
	"io"
)

type EndOfTrack struct {
	MetaEvent
}

func NewEndOfTrack(event MetaEvent, data []byte) (*EndOfTrack, error) {
	if event.eventType != 0x2f {
		return nil, fmt.Errorf("Invalid EndOfTrack event type (%02x): expected '2f'", event.eventType)
	}

	if event.length != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", event.length)
	}

	return &EndOfTrack{
		MetaEvent: event,
	}, nil
}

func (e *EndOfTrack) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                        ")

	fmt.Fprintf(w, "%02x/%-16s %s\n", e.eventType, "EndOfTrack", e.MetaEvent.Event)
}
