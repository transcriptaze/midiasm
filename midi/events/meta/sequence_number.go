package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
)

type SequenceNumber struct {
	Tag            string
	Status         types.Status
	Type           types.MetaEventType
	SequenceNumber uint16
}

func NewSequenceNumber(bytes []byte) (*SequenceNumber, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("Invalid SequenceNumber length (%d): expected '2'", len(bytes))
	}

	sequence := uint16(0)
	for _, b := range bytes {
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
