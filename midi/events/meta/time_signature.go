package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type TimeSignature struct {
	Tag                     string
	Status                  types.Status
	Type                    types.MetaEventType
	Numerator               uint8
	Denominator             uint8
	TicksPerClick           uint8
	ThirtySecondsPerQuarter uint8
}

func NewTimeSignature(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*TimeSignature, error) {
	if eventType != 0x58 {
		return nil, fmt.Errorf("Invalid TimeSignature event type (%02x): expected '58'", eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 4 {
		return nil, fmt.Errorf("Invalid TimeSignature length (%d): expected '4'", len(data))
	}

	numerator := data[0]
	ticksPerClick := data[2]
	thirtySecondsPerQuarter := data[3]

	denominator := uint8(1)
	for i := uint8(0); i < data[1]; i++ {
		denominator *= 2
	}

	return &TimeSignature{
		Tag:                     "TimeSignature",
		Status:                  status,
		Type:                    eventType,
		Numerator:               numerator,
		Denominator:             denominator,
		TicksPerClick:           ticksPerClick,
		ThirtySecondsPerQuarter: thirtySecondsPerQuarter,
	}, nil
}
