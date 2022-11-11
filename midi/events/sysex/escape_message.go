package sysex

import (
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/events"
	"github.com/transcriptaze/midiasm/midi/io"
	lib "github.com/transcriptaze/midiasm/midi/types"
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
			Status: 0xf7,
		},
		Data: data,
	}
}

func UnmarshalSysExEscapeMessage(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status lib.Status) (*SysExEscapeMessage, error) {
	if status != 0xf7 {
		return nil, fmt.Errorf("Invalid SysExEscapeMessage event type (%02x): expected 'F7'", status)
	}

	if ctx.Casio {
		return nil, fmt.Errorf("F7 is not valid for SysExEscapeMessage event in Casio mode")
	}

	data, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	event := MakeSysExEscapeMessage(tick, delta, data, r.Bytes()...)

	return &event, nil
}

// TODO encode as VLF
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

func (s *SysExEscapeMessage) UnmarshalText(bytes []byte) error {
	s.tick = 0
	s.delta = 0
	s.bytes = []byte{}
	s.tag = lib.TagSysExEscape
	s.Status = 0xf7

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)SysExEscape\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid SysExEscape event (%v)", text)
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
