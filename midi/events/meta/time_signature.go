package metaevent

import (
	"fmt"
	"regexp"
	"strconv"
)

type TimeSignature struct {
	event
	Numerator               uint8
	Denominator             uint8
	TicksPerClick           uint8
	ThirtySecondsPerQuarter uint8
}

func NewTimeSignature(tick uint64, delta uint32, bytes []byte) (*TimeSignature, error) {
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
		event: event{
			tick:  tick,
			delta: delta,
			bytes: concat([]byte{0x00, 0xff, 0x58, 0x04}, bytes),

			Tag:    "TimeSignature",
			Status: 0xff,
			Type:   0x58,
		},
		Numerator:               numerator,
		Denominator:             denominator,
		TicksPerClick:           ticksPerClick,
		ThirtySecondsPerQuarter: thirtySecondsPerQuarter,
	}, nil
}

func (t TimeSignature) String() string {
	return fmt.Sprintf("%v/%v", t.Numerator, t.Denominator)
}

func (t TimeSignature) MarshalBinary() (encoded []byte, err error) {
	d := 0
	denominator := uint8(1)
	for denominator < t.Denominator {
		denominator *= 2
		d += 1
	}

	encoded = make([]byte, 7)

	encoded[0] = byte(t.Status)
	encoded[1] = byte(t.Type)
	encoded[2] = byte(4)
	encoded[3] = byte(t.Numerator)
	encoded[4] = byte(d)
	encoded[5] = byte(t.TicksPerClick)
	encoded[6] = byte(t.ThirtySecondsPerQuarter)

	return
}

func (t *TimeSignature) UnmarshalText(bytes []byte) error {
	t.tick = 0
	t.delta = 0
	t.bytes = []byte{}
	t.Status = 0xff
	t.Tag = "TimeSignature"
	t.Type = 0x58

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)TimeSignature\s+([0-9]+)/([1-9][0-9]*),\s* ([0-9]+) ticks per click,\s*([1-9][0-8]*)/32 per quarter`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 6 {
		return fmt.Errorf("invalid TimeSignature event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if numerator, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if denominator, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if ticks, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if beats, err := strconv.ParseUint(match[5], 10, 8); err != nil {
		return err
	} else {
		t.delta = uint32(delta)
		t.Numerator = uint8(numerator)
		t.Denominator = uint8(denominator)
		t.TicksPerClick = uint8(ticks)
		t.ThirtySecondsPerQuarter = uint8(beats)
	}

	return nil
}
