package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
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

func NewTimeSignature(bytes []byte) (*TimeSignature, error) {
	if len(bytes) != 4 {
		return nil, fmt.Errorf("Invalid TimeSignature length (%d): expected '4'", len(bytes))
	}

	numerator := bytes[0]
	ticksPerClick := bytes[2]
	thirtySecondsPerQuarter := bytes[3]

	denominator := uint8(1)
	for i := uint8(0); i < bytes[1]; i++ {
		denominator *= 2
	}

	return &TimeSignature{
		Tag:                     "TimeSignature",
		Status:                  0xff,
		Type:                    0x58,
		Numerator:               numerator,
		Denominator:             denominator,
		TicksPerClick:           ticksPerClick,
		ThirtySecondsPerQuarter: thirtySecondsPerQuarter,
	}, nil
}

func (t TimeSignature) String() string {
	return fmt.Sprintf("%v/%v", t.Numerator, t.Denominator)
}
