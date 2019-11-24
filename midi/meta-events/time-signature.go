package metaevent

import (
	"fmt"
	"io"
)

type TimeSignature struct {
	MetaEvent
	numerator               uint8
	denominator             uint8
	ticksPerClick           uint8
	thirtySecondsPerQuarter uint8
}

func NewTimeSignature(event *MetaEvent, r io.ByteReader) (*TimeSignature, error) {
	if event.eventType != 0x58 {
		return nil, fmt.Errorf("Invalid TimeSignature event type (%02x): expected '58'", event.eventType)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 4 {
		return nil, fmt.Errorf("Invalid TimeSignature length (%d): expected '3'", len(data))
	}

	numerator := data[0]
	denominator := data[1]
	ticksPerClick := data[2]
	thirtySecondsPerQuarter := data[3]

	return &TimeSignature{
		MetaEvent:               *event,
		numerator:               numerator,
		denominator:             denominator,
		ticksPerClick:           ticksPerClick,
		thirtySecondsPerQuarter: thirtySecondsPerQuarter,
	}, nil
}

func (e *TimeSignature) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                            ")

	fmt.Fprintf(w, "%02x/%-16s %s numerator:%d denominator:%d ticks/click:%d 1/32-per-quarter:%d\n", e.eventType, "TimeSignature", e.MetaEvent.Event, e.numerator, e.denominator, e.ticksPerClick, e.thirtySecondsPerQuarter)
}
