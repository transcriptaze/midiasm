package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
)

type Tempo struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Tempo  uint32
}

func NewTempo(r events.EventReader, status types.Status, eventType types.MetaEventType) (*Tempo, error) {
	if eventType != 0x51 {
		return nil, fmt.Errorf("Invalid Tempo event type (%02x): expected '51'", eventType)
	}

	data, err := r.ReadVLF()
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
		Tag:    "Tempo",
		Status: status,
		Type:   eventType,
		Tempo:  tempo,
	}, nil
}
