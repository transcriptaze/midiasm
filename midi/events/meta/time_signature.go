package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type TimeSignature struct {
	event
	Numerator               uint8
	Denominator             uint8
	TicksPerClick           uint8
	ThirtySecondsPerQuarter uint8
}

func MakeTimeSignature(tick uint64, delta lib.Delta, numerator, denominator, ticksPerClick, thirtySecondsPerQuarter uint8, bytes ...byte) TimeSignature {
	d := 0
	dd := uint8(1)
	for dd < denominator {
		dd *= 2
		d += 1
	}

	return TimeSignature{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagTimeSignature,
			Status: 0xff,
			Type:   lib.TypeTimeSignature,
		},
		Numerator:               numerator,
		Denominator:             denominator,
		TicksPerClick:           ticksPerClick,
		ThirtySecondsPerQuarter: thirtySecondsPerQuarter,
	}
}

func (e *TimeSignature) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	if len(data) != 4 {
		return fmt.Errorf("Invalid TimeSignature length (%d): expected '4'", len(data))
	}

	numerator := data[0]
	ticksPerClick := data[2]
	thirtySecondsPerQuarter := data[3]

	denominator := uint8(1)
	for i := uint8(0); i < data[1]; i++ {
		denominator *= 2
	}

	*e = MakeTimeSignature(tick, delta, numerator, denominator, ticksPerClick, thirtySecondsPerQuarter, bytes...)

	return nil
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

func (e *TimeSignature) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagTimeSignature, remaining[0])
	} else if !equals(remaining[1], lib.TypeTimeSignature) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagTimeSignature, remaining[1])
	} else if data, err := vlf(remaining[2:]); err != nil {
		return err
	} else if len(data) < 4 {
		return fmt.Errorf("Invalid TimeSignature data")
	} else {
		numerator := data[0]
		ticksPerClick := data[2]
		thirtySecondsPerQuarter := data[3]

		denominator := uint8(1)
		for i := uint8(0); i < data[1]; i++ {
			denominator *= 2
		}

		*e = MakeTimeSignature(0, delta, numerator, denominator, ticksPerClick, thirtySecondsPerQuarter, bytes...)
	}

	return nil
}

func (e *TimeSignature) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)TimeSignature\s+([0-9]+)/([1-9][0-9]*),\s* ([0-9]+) ticks per click,\s*([1-9][0-8]*)/32 per quarter`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 6 {
		return fmt.Errorf("invalid TimeSignature event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
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
		*e = MakeTimeSignature(0, delta, uint8(numerator), uint8(denominator), uint8(ticks), uint8(beats), []byte{}...)
	}

	return nil
}

func (e TimeSignature) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag                     string    `json:"tag"`
		Delta                   lib.Delta `json:"delta"`
		Status                  byte      `json:"status"`
		Type                    byte      `json:"type"`
		Numerator               uint8     `json:"numerator"`
		Denominator             uint8     `json:"denominator"`
		TicksPerClick           uint8     `json:"ticks-per-click"`
		ThirtySecondsPerQuarter uint8     `json:"thirty-seconds-per-quarter"`
	}{
		Tag:                     fmt.Sprintf("%v", e.tag),
		Delta:                   e.delta,
		Status:                  byte(e.Status),
		Type:                    byte(e.Type),
		Numerator:               e.Numerator,
		Denominator:             e.Denominator,
		TicksPerClick:           e.TicksPerClick,
		ThirtySecondsPerQuarter: e.ThirtySecondsPerQuarter,
	}

	return json.Marshal(t)
}

func (e *TimeSignature) UnmarshalJSON(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.tag = lib.TagTimeSignature
	e.Type = lib.TypeTimeSignature

	t := struct {
		Tag                     string    `json:"tag"`
		Delta                   lib.Delta `json:"delta"`
		Numerator               uint8     `json:"numerator"`
		Denominator             uint8     `json:"denominator"`
		TicksPerClick           uint8     `json:"ticks-per-click"`
		ThirtySecondsPerQuarter uint8     `json:"thirty-seconds-per-quarter"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if t.Tag != fmt.Sprintf("%v", lib.TagTimeSignature) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeTimeSignature(0, t.Delta, t.Numerator, t.Denominator, t.TicksPerClick, t.ThirtySecondsPerQuarter, []byte{}...)
	}

	return nil
}
