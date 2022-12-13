package metaevent

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type SequenceNumber struct {
	event
	SequenceNumber uint16
}

func MakeSequenceNumber(tick uint64, delta lib.Delta, sequence uint16, bytes ...byte) SequenceNumber {
	return SequenceNumber{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagSequenceNumber,
			Status: 0xff,
			Type:   lib.TypeSequenceNumber,
		},
		SequenceNumber: sequence,
	}
}

func (e *SequenceNumber) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	if len(data) != 2 {
		return fmt.Errorf("Invalid SequenceNumber length (%d): expected '2'", len(data))
	}

	sequence := binary.BigEndian.Uint16(data)
	event := MakeSequenceNumber(tick, delta, sequence, bytes...)

	*e = event

	return nil
}

func (s SequenceNumber) MarshalBinary() (encoded []byte, err error) {
	encoded = binary.BigEndian.AppendUint16([]byte{
		byte(s.Status),
		byte(s.Type),
		byte(2),
	}, s.SequenceNumber)

	return
}

func (e *SequenceNumber) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagSequenceNumber, remaining[0])
	} else if !equals(remaining[1], lib.TypeSequenceNumber) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagSequenceNumber, remaining[1])
	} else if v, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeSequenceNumber(0, delta, binary.BigEndian.Uint16(v), bytes...)
	}

	return nil
}

func (e *SequenceNumber) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SequenceNumber\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid SequenceNumber event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if sequence, err := strconv.ParseUint(match[2], 10, 16); err != nil {
		return err
	} else {
		*e = MakeSequenceNumber(0, delta, uint16(sequence), []byte{}...)
	}

	return nil
}

func (e SequenceNumber) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag            string    `json:"tag"`
		Delta          lib.Delta `json:"delta"`
		Status         byte      `json:"status"`
		Type           byte      `json:"type"`
		SequenceNumber uint16    `json:"sequence-number"`
	}{
		Tag:            fmt.Sprintf("%v", e.tag),
		Delta:          e.delta,
		Status:         byte(e.Status),
		Type:           byte(e.Type),
		SequenceNumber: e.SequenceNumber,
	}

	return json.Marshal(t)
}

func (e *SequenceNumber) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag            string    `json:"tag"`
		Delta          lib.Delta `json:"delta"`
		SequenceNumber uint16    `json:"sequence-number"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagSequenceNumber) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeSequenceNumber(0, t.Delta, t.SequenceNumber, []byte{}...)
	}

	return nil
}
