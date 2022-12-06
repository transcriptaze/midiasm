package sysex

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type SysExContinuationMessage struct {
	event
	Data lib.Hex
	End  bool
}

func MakeSysExContinuationMessage(tick uint64, delta uint32, data lib.Hex, bytes ...byte) SysExContinuationMessage {
	return SysExContinuationMessage{
		event: event{
			tick:   tick,
			delta:  lib.Delta(delta),
			bytes:  bytes,
			tag:    lib.TagSysExContinuation,
			Status: 0xf7,
		},
		Data: data,
		End:  false,
	}
}

func MakeSysExContinuationEndMessage(tick uint64, delta uint32, data lib.Hex, bytes ...byte) SysExContinuationMessage {
	return SysExContinuationMessage{
		event: event{
			tick:   tick,
			delta:  lib.Delta(delta),
			bytes:  bytes,
			tag:    lib.TagSysExContinuation,
			Status: 0xf7,
		},
		Data: data,
		End:  true,
	}
}

func (e *SysExContinuationMessage) unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status != 0xf7 {
		return fmt.Errorf("Invalid SysExContinuationMessage event type (%02x): expected 'F7'", status)
	}

	d := data
	if len(data) > 0 {
		terminator := data[len(d)-1]
		if terminator == 0xf7 {
			d = data[:len(d)-1]
			ctx.Casio = false
		}
	}

	*e = MakeSysExContinuationMessage(tick, delta, d, bytes...)

	return nil
}

func (s SysExContinuationMessage) MarshalBinary() ([]byte, error) {
	status := byte(s.Status)

	vlf := []byte{}
	vlf = append(vlf, s.Data...)

	if s.End {
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

func (s *SysExContinuationMessage) UnmarshalText(bytes []byte) error {
	s.tick = 0
	s.delta = 0
	s.bytes = []byte{}
	s.tag = lib.TagSysExContinuation
	s.Status = 0xf7

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SysExContinuation\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid SysExContinuation event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if data, err := lib.ParseHex(match[2]); err != nil {
		return err
	} else {
		s.delta = delta
		s.Data = data
	}

	return nil
}

func (e SysExContinuationMessage) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Data   lib.Hex   `json:"data"`
		End    bool      `json:"end"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Data:   e.Data,
		End:    e.End,
	}

	return json.Marshal(t)
}

func (e *SysExContinuationMessage) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Data  lib.Hex   `json:"data"`
		End   bool      `json:"end"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagSysExContinuation) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.tag = lib.TagSysExContinuation
		e.Status = 0xF7
		e.Data = t.Data
		e.End = t.End
	}

	return nil
}
