package metaevent

import (
	"fmt"
	"io"
)

type SequenceNumber struct {
	Tag string
	MetaEvent
	SequenceNumber uint16
}

func NewSequenceNumber(event *MetaEvent, r io.ByteReader) (*SequenceNumber, error) {
	if event.Type != 0x00 {
		return nil, fmt.Errorf("Invalid SequenceNumber event type (%02x): expected '00'", event.Type)
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
		MetaEvent:      *event,
		SequenceNumber: sequence,
	}, nil
}
