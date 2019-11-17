package metaevent

import (
	"fmt"
	"io"
)

type Tempo struct {
	MetaEvent
	Tempo uint32
}

func NewTempo(event MetaEvent, data []byte) (*Tempo, error) {
	if event.status != 0xff {
		return nil, fmt.Errorf("Invalid Tempo status (%02x): expected 'ff'", event.status)
	}

	if event.eventType != 0x51 {
		return nil, fmt.Errorf("Invalid Tempo event type (%02x): expected '51'", event.eventType)
	}

	if event.length != 3 {
		return nil, fmt.Errorf("Invalid Tempo length (%d): expected '3'", event.length)
	}

	tempo := uint32(0)
	for i := 0; i < int(event.length); i++ {
		tempo <<= 8
		tempo += uint32(data[i])
	}

	return &Tempo{
		MetaEvent: event,
		Tempo:     tempo,
	}, nil
}

func (e *Tempo) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                               ")

	fmt.Fprintf(w, "%02x/%-16s delta:%-10d tempo:%v\n", e.eventType, "Tempo", e.delta, e.Tempo)
}
