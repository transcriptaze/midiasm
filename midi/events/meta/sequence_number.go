package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type SequenceNumber struct {
	Tag            string
	Status         types.Status
	Type           types.MetaEventType
	SequenceNumber uint16
}

func NewSequenceNumber(r io.ByteReader) (*SequenceNumber, error) {
	data, err := events.VLF(r)
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
		Status:         0xff,
		Type:           0x00,
		SequenceNumber: sequence,
	}, nil
}
