package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Tempo struct {
	Tag string
	*events.Event
	Type  types.MetaEventType
	Tempo uint32
}

func NewTempo(event *events.Event, eventType types.MetaEventType, r io.ByteReader) (*Tempo, error) {
	if eventType != 0x51 {
		return nil, fmt.Errorf("Invalid Tempo event type (%02x): expected '51'", eventType)
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
		Tag:   "Tempo",
		Event: event,
		Type:  eventType,
		Tempo: tempo,
	}, nil
}
