package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type InstrumentName struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Name   string
}

func NewInstrumentName(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*InstrumentName, error) {
	if eventType != 0x04 {
		return nil, fmt.Errorf("Invalid InstrumentName event type (%02x): expected '04'", eventType)
	}

	name, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &InstrumentName{
		Tag:    "InstrumentName",
		Status: status,
		Type:   eventType,
		Name:   string(name),
	}, nil
}
