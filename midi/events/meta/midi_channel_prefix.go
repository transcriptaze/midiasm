package metaevent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type MIDIChannelPrefix struct {
	event
	Channel uint8
}

// func NewMIDIChannelPrefix(bytes []byte) (*MIDIChannelPrefix, error) {
// 	if len(bytes) != 1 {
// 		return nil, fmt.Errorf("Invalid MIDIChannelPrefix length (%d): expected '1'", len(bytes))
// 	}

// 	channel := int8(bytes[0])
// 	if channel < 0 || channel > 15 {
// 		return nil, fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
// 	}

// 	return &MIDIChannelPrefix{
// 		Tag:     "MIDIChannelPrefix",
// 		Status:  0xff,
// 		Type:    0x20,
// 		Channel: channel,
// 	}, nil
// }

func MakeMIDIChannelPrefix(tick uint64, delta uint32, channel uint8) MIDIChannelPrefix {
	return MIDIChannelPrefix{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x20, 0x01, channel}),
			tag:    types.TagMIDIChannelPrefix,
			Status: 0xff,
			Type:   types.TypeMIDIChannelPrefix,
		},
		Channel: channel,
	}
}

func UnmarshalMIDIChannelPrefix(tick uint64, delta uint32, bytes []byte) (*MIDIChannelPrefix, error) {
	if len(bytes) != 1 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix length (%d): expected '1'", len(bytes))
	}

	channel := bytes[0]
	event := MakeMIDIChannelPrefix(tick, delta, channel)

	return &event, nil
}

func (m MIDIChannelPrefix) MarshalBinary() (encoded []byte, err error) {
	return []byte{
		byte(m.Status),
		byte(m.Type),
		1,
		m.Channel,
	}, nil
}

func (m *MIDIChannelPrefix) UnmarshalText(bytes []byte) error {
	m.tick = 0
	m.delta = 0
	m.bytes = []byte{}
	m.tag = types.TagMIDIChannelPrefix
	m.Status = 0xff
	m.Type = types.TypeMIDIChannelPrefix

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)MIDIChannelPrefix\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid MIDIChannelPrefix event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if channel < 0 || channel > 15 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	} else {
		m.delta = uint32(delta)
		m.Channel = uint8(channel)
	}

	return nil
}
