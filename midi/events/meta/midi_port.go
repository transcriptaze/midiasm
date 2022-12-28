package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

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

func (e *MIDIPort) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	if len(data) != 1 {
		return fmt.Errorf("Invalid MIDIPort length (%d): expected '1'", len(data))
	}

	port := data[0]

	if port > 127 {
		return fmt.Errorf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port)
	}

	*e = MakeMIDIPort(tick, delta, port, bytes...)

	return nil
}

func (m MIDIPort) MarshalBinary() (encoded []byte, err error) {
	return []byte{
		byte(m.Status),
		byte(m.Type),
		1,
		m.Port,
	}, nil
}

func (e *MIDIPort) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagMIDIPort, remaining[0])
	} else if !equals(remaining[1], lib.TypeMIDIPort) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagMIDIPort, remaining[1])
	} else if data, err := vlf(remaining[2:]); err != nil {
		return err
	} else if len(data) < 1 {
		return fmt.Errorf("Invalid MIDIPort channel data")
	} else if port := data[0]; port > 127 {
		return fmt.Errorf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port)
	} else {
		*e = MakeMIDIPort(0, delta, port, bytes...)
	}

	return nil
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
	} else if port > 127 {
		return fmt.Errorf("Invalid MIDIPort channel (%d): expected a value in the interval [0..127]", port)
	} else {
		*e = MakeMIDIPort(0, delta, uint8(port), []byte{}...)
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
	} else if t.Port > 127 {
		return fmt.Errorf("Invalid MIDIPort channel (%d): expected a value in the interval [0..127]", t.Port)
	} else {
		*e = MakeMIDIPort(0, t.Delta, t.Port, []byte{}...)
	}

	return nil
}
