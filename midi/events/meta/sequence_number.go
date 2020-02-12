package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SequenceNumber struct {
	Tag string
	*events.Event
	Type           types.MetaEventType
	SequenceNumber uint16
}

func NewSequenceNumber(event *events.Event, eventType types.MetaEventType, r io.ByteReader) (*SequenceNumber, error) {
	if eventType != 0x00 {
		return nil, fmt.Errorf("Invalid SequenceNumber event type (%02x): expected '00'", eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 2 {
		return nil, fmt.Errorf("Invalid SequenceNumber length (%d): expected '2'", len(data))
	}

	sequence := uint16(0)
	for _, b := range data {
		sequence <<= 8
		sequence += uint16(b)
	}

	return &SequenceNumber{
		Tag:            "SequenceNumber",
		Event:          event,
		Type:           eventType,
		SequenceNumber: sequence,
	}, nil
}
