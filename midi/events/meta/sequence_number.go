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
			tick:  tick,
			delta: delta,
			// bytes:  binary.BigEndian.AppendUint16([]byte{0x00, 0xff, 0x00, 0x02}, sequence),
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
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagSequenceNumber
		e.Status = 0xff
		e.Type = lib.TypeSequenceNumber
		e.SequenceNumber = uint16(sequence)
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
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagSequenceNumber
		e.Type = lib.TypeSequenceNumber
		e.SequenceNumber = t.SequenceNumber
	}

	return nil
}
