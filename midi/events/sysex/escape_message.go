package sysex

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type SysExEscapeMessage struct {
	event
	Data lib.Hex
}

func MakeSysExEscapeMessage(tick uint64, delta uint32, data lib.Hex, bytes ...byte) SysExEscapeMessage {
	return SysExEscapeMessage{
		event: event{
			tick:   tick,
			delta:  lib.Delta(delta),
			bytes:  bytes,
			tag:    lib.TagSysExEscape,
			Status: 0xF7,
		},
		Data: data,
	}
}

func (e *SysExEscapeMessage) unmarshal(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status != 0xF7 {
		return fmt.Errorf("Invalid SysExEscapeMessage event type (%02x): expected 'F7'", status)
	}

	*e = MakeSysExEscapeMessage(tick, delta, data, bytes...)

	return nil
}

func (s SysExEscapeMessage) MarshalBinary() ([]byte, error) {
	status := byte(s.Status)

	vlf := []byte{}
	vlf = append(vlf, s.Data...)

	if data, err := lib.VLF(vlf).MarshalBinary(); err != nil {
		return nil, err
	} else {
		encoded := []byte{}
		encoded = append(encoded, status)
		encoded = append(encoded, data...)

		return encoded, nil
	}
}

func (e *SysExEscapeMessage) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeSysExEscapeMessage.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagSysExEscape, remaining[0])
	} else if data, err := vlf(remaining[1:]); err != nil {
		return err
	} else {
		*e = MakeSysExEscapeMessage(0, uint32(delta), data, bytes...)
	}

	return nil
}

func (e *SysExEscapeMessage) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SysExEscape\s+(.*)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid SysExEscape event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if data, err := lib.ParseHex(match[2]); err != nil {
		return err
	} else {
		*e = MakeSysExEscapeMessage(0, uint32(delta), data, []byte{}...)
	}

	return nil
}

func (e SysExEscapeMessage) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Data   lib.Hex   `json:"data"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Data:   e.Data,
	}

	return json.Marshal(t)
}

func (e *SysExEscapeMessage) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Data  lib.Hex   `json:"data"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagSysExEscape) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeSysExEscapeMessage(0, uint32(t.Delta), t.Data, []byte{}...)
	}

	return nil
}
