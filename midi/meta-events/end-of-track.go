package metaevent

import (
	"fmt"
	"io"
)

type EndOfTrack struct {
	MetaEvent
}

func NewEndOfTrack(event MetaEvent, data []byte) (*EndOfTrack, error) {
	if event.status != 0xff {
		return nil, fmt.Errorf("Invalid EndOfTrack status (%02x): expected 'ff'", event.status)
	}

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

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d\n", e.eventType, "EndOfTrack", e.delta)
}
