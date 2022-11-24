package metaevent

import (
	"encoding/hex"
	"fmt"
	"regexp"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type SequencerSpecificEvent struct {
	event
	Manufacturer lib.Manufacturer
	Data         lib.Hex
}

func MakeSequencerSpecificEvent(tick uint64, delta lib.Delta, manufacturer lib.Manufacturer, data []byte) SequencerSpecificEvent {
	return SequencerSpecificEvent{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  concat([]byte{0x00, 0xff, 0x7f, byte(len(manufacturer.ID) + len(data))}, manufacturer.ID, data),
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
		},
		Manufacturer: manufacturer,
		Data:         data,
	}
}

func UnmarshalSequencerSpecificEvent(tick uint64, delta lib.Delta, bytes []byte) (*SequencerSpecificEvent, error) {
	id := bytes[0:1]
	data := bytes[1:]

	if bytes[0] == 0x00 {
		id = bytes[0:3]
		data = bytes[3:]
	}

	event := MakeSequencerSpecificEvent(tick, delta, lib.LookupManufacturer(id), data)

	return &event, nil
}

func (s SequencerSpecificEvent) MarshalBinary() (encoded []byte, err error) {
	return concat(
		[]byte{
			byte(s.Status),
			byte(s.Type),
			byte(len(s.Manufacturer.ID) + len(s.Data))},
		s.Manufacturer.ID,
		s.Data), nil
}

func (s *SequencerSpecificEvent) UnmarshalText(bytes []byte) error {
	s.tick = 0
	s.delta = 0
	s.bytes = []byte{}
	s.tag = lib.TagSequencerSpecificEvent
	s.Status = 0xff
	s.Type = lib.TypeSequencerSpecificEvent

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SequencerSpecificEvent\s+(.*?),(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("Invalid SequencerSpecificEvent (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if manufacturer, err := lib.FindManufacturer(match[2]); err != nil {
		return err
	} else {
		data := []byte{}
		if len(match) > 3 {
			s := regexp.MustCompile(`\s+`).ReplaceAllString(match[3], "")
			if d, err := hex.DecodeString(s); err != nil {
				return err
			} else {
				data = d
			}
		}

		s.delta = delta
		s.Manufacturer = manufacturer
		s.Data = data
	}

	return nil
}
