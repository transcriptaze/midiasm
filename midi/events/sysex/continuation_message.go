package sysex

import (
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	lib "github.com/transcriptaze/midiasm/midi/types"
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

func UnmarshalSysExContinuationMessage(ctx *context.Context, tick uint64, delta uint32, status lib.Status, bytes []byte, src ...byte) (*SysExContinuationMessage, error) {
	if status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExContinuationMessage event type (%02x): expected 'F7'", status)
	}

	data := bytes
	if len(bytes) > 0 {
		terminator := bytes[len(data)-1]
		if terminator == 0xf7 {
			data = bytes[:len(data)-1]
			ctx.Casio = false
		}
	}

	event := MakeSysExContinuationMessage(tick, delta, data, src...)

	return &event, nil
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
