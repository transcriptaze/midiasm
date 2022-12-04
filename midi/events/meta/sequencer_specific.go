package metaevent

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type SequencerSpecificEvent struct {
	event
	Manufacturer lib.Manufacturer
	Data         lib.Hex
}

func MakeSequencerSpecificEvent(tick uint64, delta lib.Delta, manufacturer lib.Manufacturer, data []byte, bytes ...byte) SequencerSpecificEvent {
	return SequencerSpecificEvent{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagSequencerSpecificEvent,
			Status: 0xff,
			Type:   lib.TypeSequencerSpecificEvent,
		},
		Manufacturer: manufacturer,
		Data:         data,
	}
}

func UnmarshalSequencerSpecificEvent(ctx *context.Context, tick uint64, delta lib.Delta, data []byte, bytes ...byte) (*SequencerSpecificEvent, error) {
	id := data[0:1]
	d := data[1:]

	if data[0] == 0x00 {
		id = data[0:3]
		d = data[3:]
	}

	event := MakeSequencerSpecificEvent(tick, delta, lib.LookupManufacturer(id), d, bytes...)

	return &event, nil
}

// FIXME encode as VLF
func (e SequencerSpecificEvent) MarshalBinary() (encoded []byte, err error) {
	b := []byte{
		byte(e.Status),
		byte(e.Type),
		byte(len(e.Manufacturer.ID) + len(e.Data))}

	b = append(b, e.Manufacturer.ID...)
	b = append(b, e.Data...)

	return b, nil
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

func (e SequencerSpecificEvent) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag          string           `json:"tag"`
		Delta        lib.Delta        `json:"delta"`
		Status       byte             `json:"status"`
		Type         byte             `json:"type"`
		Manufacturer lib.Manufacturer `json:"manufacturer"`
		Data         lib.Hex          `json:"data"`
	}{
		Tag:          fmt.Sprintf("%v", e.tag),
		Delta:        e.delta,
		Status:       byte(e.Status),
		Type:         byte(e.Type),
		Manufacturer: e.Manufacturer,
		Data:         e.Data,
	}

	return json.Marshal(t)
}

func (e *SequencerSpecificEvent) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag          string           `json:"tag"`
		Delta        lib.Delta        `json:"delta"`
		Manufacturer lib.Manufacturer `json:"manufacturer"`
		Data         lib.Hex          `json:"data"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagSequencerSpecificEvent) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagSequencerSpecificEvent
		e.Type = lib.TypeSequencerSpecificEvent
		e.Manufacturer = t.Manufacturer
		e.Data = t.Data
	}

	return nil
}
