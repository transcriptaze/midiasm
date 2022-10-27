package metaevent

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type SequenceNumber struct {
	event
	SequenceNumber uint16
}

func NewSequenceNumber(tick uint64, delta uint32, sequence uint16) (*SequenceNumber, error) {
	return &SequenceNumber{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  binary.BigEndian.AppendUint16([]byte{0x00, 0xff, 0x00, 0x02}, sequence),
			tag:    types.TagSequenceNumber,
			Status: 0xff,
			Type:   0x00,
		},
		SequenceNumber: sequence,
	}, nil
}

func UnmarshalSequenceNumber(tick uint64, delta uint32, bytes []byte) (*SequenceNumber, error) {
	if len(bytes) != 2 {
		return nil, fmt.Errorf("Invalid SequenceNumber length (%d): expected '2'", len(bytes))
	}

	sequence := binary.BigEndian.Uint16(bytes)

	return &SequenceNumber{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  concat([]byte{0x00, 0xff, 0x00, 0x02}, bytes),
			tag:    types.TagSequenceNumber,
			Status: 0xff,
			Type:   0x00,
		},
		SequenceNumber: sequence,
	}, nil
}

func (s SequenceNumber) MarshalBinary() (encoded []byte, err error) {
	encoded = binary.BigEndian.AppendUint16([]byte{
		byte(s.Status),
		byte(s.Type),
		byte(2),
	}, s.SequenceNumber)

	return
}

func (s *SequenceNumber) UnmarshalText(bytes []byte) error {
	s.tick = 0
	s.delta = 0
	s.bytes = []byte{}
	s.tag = types.TagSequenceNumber
	s.Status = 0xff
	s.Type = 0x00

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SequenceNumber\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid SequenceNumber event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if sequence, err := strconv.ParseUint(match[2], 10, 16); err != nil {
		return err
	} else {
		s.delta = uint32(delta)
		s.SequenceNumber = uint16(sequence)
	}

	return nil
}
