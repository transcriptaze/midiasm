package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type MIDIPort struct {
	event
	Port uint8
}

func MakeMIDIPort(tick uint64, delta uint32, port uint8) MIDIPort {
	if port > 127 {
		panic(fmt.Sprintf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port))
	}

	return MIDIPort{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x21, 0x01, port}),
			tag:    types.TagMIDIPort,
			Status: 0xff,
			Type:   types.TypeMIDIPort,
		},
		Port: port,
	}
}

func UnmarshalMIDIPort(tick uint64, delta uint32, bytes []byte) (*MIDIPort, error) {
	if len(bytes) != 1 {
		return nil, fmt.Errorf("Invalid MIDIPort length (%d): expected '1'", len(bytes))
	}

	if port := bytes[0]; port > 127 {
		return nil, fmt.Errorf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port)
	} else {
		event := MakeMIDIPort(tick, delta, port)

		return &event, nil
	}
}

func (m MIDIPort) MarshalBinary() (encoded []byte, err error) {
	return []byte{
		byte(m.Status),
		byte(m.Type),
		1,
		m.Port,
	}, nil
}

func (m *MIDIPort) UnmarshalText(bytes []byte) error {
	m.tick = 0
	m.delta = 0
	m.bytes = []byte{}
	m.tag = types.TagMIDIPort
	m.Status = 0xff
	m.Type = types.TypeMIDIPort

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)MIDIPort\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("Invalid MIDIPort event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if port, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if port > 1278 {
		return fmt.Errorf("Invalid MIDIPort channel (%d): expected a value in the interval [0..127]", port)
	} else {
		m.delta = uint32(delta)
		m.Port = uint8(port)
	}

	return nil
}
