package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type InstrumentName struct {
	Tag string
	MetaEvent
	Type types.MetaEventType
	Name string
}

func NewInstrumentName(event *MetaEvent, eventType types.MetaEventType, r io.ByteReader) (*InstrumentName, error) {
	if eventType != 0x04 {
		return nil, fmt.Errorf("Invalid InstrumentName event type (%02x): expected '04'", eventType)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &InstrumentName{
		Tag:       "InstrumentName",
		MetaEvent: *event,
		Type:      eventType,
		Name:      string(name),
	}, nil
}
