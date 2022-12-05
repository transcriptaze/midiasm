package sysex

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
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

func UnmarshalSysExMessage(ctx *context.Context, tick uint64, delta uint32, status lib.Status, bytes []byte, src ...byte) (*SysExMessage, error) {
	if status != 0xf0 {
		return nil, fmt.Errorf("Invalid SysExMessage event type (%02x): expected 'F0'", status)
	}

	if len(bytes) == 0 {
		return nil, fmt.Errorf("Invalid SysExMessage event data (0 bytes")
	}

	id := bytes[0:1]
	manufacturer := lib.LookupManufacturer(id)
	data := bytes[1:]

	if data[len(data)-1] != 0xf7 {
		ctx.Casio = true
		event := MakeSysExMessage(tick, delta, manufacturer, data, src...)
		return &event, nil
	} else {
		data = data[:len(data)-1]
		event := MakeSysExSingleMessage(tick, delta, manufacturer, data, src...)
		return &event, nil
	}
}

func (s SysExMessage) MarshalBinary() ([]byte, error) {
	status := byte(s.Status)

	vlf := []byte{}
	vlf = append(vlf, s.Manufacturer.ID...)
	vlf = append(vlf, s.Data...)

	if s.Single {
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
