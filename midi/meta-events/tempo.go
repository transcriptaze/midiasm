package metaevent

import (
	"fmt"
	"io"
)

type Tempo struct {
	MetaEvent
	Tempo uint32
}

func NewTempo(event *MetaEvent, r io.ByteReader) (*Tempo, error) {
	if event.eventType != 0x51 {
		return nil, fmt.Errorf("Invalid Tempo event type (%02x): expected '51'", event.eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 3 {
		return nil, fmt.Errorf("Invalid Tempo length (%d): expected '3'", len(data))
	}

	tempo := uint32(0)
	for _, b := range data {
		tempo <<= 8
		tempo += uint32(b)
	}

	return &Tempo{
		MetaEvent: *event,
		Tempo:     tempo,
	}, nil
}

func (e *Tempo) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                               ")

	fmt.Fprintf(w, "%02x/%-16s %s tempo:%v\n", e.eventType, "Tempo", e.MetaEvent.Event, e.Tempo)
}
