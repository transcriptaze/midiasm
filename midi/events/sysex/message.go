package sysex

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type SysExMessage struct {
	event
	Manufacturer lib.Manufacturer
	Data         lib.Hex
	Single       bool
}

func MakeSysExMessage(tick uint64, delta uint32, manufacturer lib.Manufacturer, data lib.Hex, bytes ...byte) SysExMessage {
	return SysExMessage{
		event: event{
			tick:   tick,
			delta:  lib.Delta(delta),
			bytes:  bytes,
			tag:    lib.TagSysExMessage,
			Status: 0xF0,
		},
		Manufacturer: manufacturer,
		Data:         data,
		Single:       false,
	}
}

func MakeSysExSingleMessage(tick uint64, delta uint32, manufacturer lib.Manufacturer, data lib.Hex, bytes ...byte) SysExMessage {
	return SysExMessage{
		event: event{
			tick:   tick,
			delta:  lib.Delta(delta),
			bytes:  bytes,
			tag:    lib.TagSysExMessage,
			Status: 0xf0,
		},
		Manufacturer: manufacturer,
		Data:         data,
		Single:       true,
	}
}

func (e *SysExMessage) unmarshal(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status != 0xf0 {
		return fmt.Errorf("Invalid SysExMessage event type (%02x): expected 'F0'", status)
	}

	if len(data) == 0 {
		return fmt.Errorf("Invalid SysExMessage event data (0 bytes")
	}

	id := data[0:1]
	manufacturer := lib.LookupManufacturer(id)
	d := data[1:]

	if d[len(d)-1] != 0xf7 {
		*e = MakeSysExMessage(tick, delta, manufacturer, d, bytes...)
	} else {
		*e = MakeSysExSingleMessage(tick, delta, manufacturer, d[:len(d)-1], bytes...)
	}

	return nil
}

func (e SysExMessage) MarshalBinary() ([]byte, error) {
	status := byte(e.Status)

	vlf := []byte{}
	vlf = append(vlf, e.Manufacturer.ID...)
	vlf = append(vlf, e.Data...)

	if e.Single {
		vlf = append(vlf, 0xf7)
	}

	if data, err := lib.VLF(vlf).MarshalBinary(); err != nil {
		return nil, err
	} else {
		encoded := []byte{}
		encoded = append(encoded, status)
		encoded = append(encoded, data...)

		return encoded, nil
	}
}

func (e *SysExMessage) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeSysExMessage.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagSysExMessage, remaining[0])
	} else if data, err := vlf(remaining[1:]); err != nil {
		return err
	} else {
		id := data[0:1]
		manufacturer := lib.LookupManufacturer(id)
		d := data[1:]

		if d[len(d)-1] != 0xf7 {
			// ctx.Casio = true
			*e = MakeSysExMessage(0, uint32(delta), manufacturer, d, bytes...)
		} else {
			*e = MakeSysExSingleMessage(0, uint32(delta), manufacturer, d[:len(d)-1], bytes...)
		}

	}

	return nil
}

func (s *SysExMessage) UnmarshalText(bytes []byte) error {
	s.tick = 0
	s.delta = 0
	s.bytes = []byte{}
	s.tag = lib.TagSysExMessage
	s.Status = 0xf0

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SysExMessage\s+(.*?),(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid SysExMessage event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if manufacturer, err := lib.FindManufacturer(match[2]); err != nil {
		return err
	} else if data, err := lib.ParseHex(match[3]); err != nil {
		return err
	} else {
		s.delta = delta
		s.Manufacturer = manufacturer
		s.Data = data
	}

	return nil
}

func (e SysExMessage) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag          string           `json:"tag"`
		Delta        lib.Delta        `json:"delta"`
		Status       byte             `json:"status"`
		Manufacturer lib.Manufacturer `json:"manufacturer"`
		Data         lib.Hex          `json:"data"`
		Single       bool             `json:"single"`
	}{
		Tag:          fmt.Sprintf("%v", e.tag),
		Delta:        e.delta,
		Status:       byte(e.Status),
		Manufacturer: e.Manufacturer,
		Data:         e.Data,
		Single:       e.Single,
	}

	return json.Marshal(t)
}

func (e *SysExMessage) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag          string           `json:"tag"`
		Delta        lib.Delta        `json:"delta"`
		Manufacturer lib.Manufacturer `json:"manufacturer"`
		Data         lib.Hex          `json:"data"`
		Single       bool             `json:"single"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagSysExMessage) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.tag = lib.TagSysExMessage
		e.Status = 0xF0
		e.Manufacturer = t.Manufacturer
		e.Data = t.Data
		e.Single = t.Single
	}

	return nil
}
