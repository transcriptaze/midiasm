package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type MIDIPort struct {
	event
	Port uint8
}

func MakeMIDIPort(tick uint64, delta lib.Delta, port uint8, bytes ...byte) MIDIPort {
	if port > 127 {
		panic(fmt.Sprintf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port))
	}

	return MIDIPort{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagMIDIPort,
			Status: 0xff,
			Type:   lib.TypeMIDIPort,
		},
		Port: port,
	}
}

func UnmarshalMIDIPort(ctx *context.Context, tick uint64, delta lib.Delta, data []byte, bytes ...byte) (*MIDIPort, error) {
	if len(data) != 1 {
		return nil, fmt.Errorf("Invalid MIDIPort length (%d): expected '1'", len(data))
	}

	port := data[0]

	if port > 127 {
		return nil, fmt.Errorf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port)
	}

	event := MakeMIDIPort(tick, delta, port, bytes...)

	return &event, nil
}

func (m MIDIPort) MarshalBinary() (encoded []byte, err error) {
	return []byte{
		byte(m.Status),
		byte(m.Type),
		1,
		m.Port,
	}, nil
}

func (e *MIDIPort) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)MIDIPort\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("Invalid MIDIPort event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if port, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if port > 1278 {
		return fmt.Errorf("Invalid MIDIPort channel (%d): expected a value in the interval [0..127]", port)
	} else {
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagMIDIPort
		e.Status = 0xff
		e.Type = lib.TypeMIDIPort
		e.Port = uint8(port)
	}

	return nil
}

func (e MIDIPort) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Type   byte      `json:"type"`
		Port   uint8     `json:"port"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
		Port:   e.Port,
	}

	return json.Marshal(t)
}

func (e *MIDIPort) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Port  uint8     `json:"port"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagMIDIPort) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagMIDIPort
		e.Type = lib.TypeMIDIPort
		e.Port = t.Port
	}

	return nil
}
