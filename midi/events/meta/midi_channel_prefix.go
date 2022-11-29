package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type MIDIChannelPrefix struct {
	event
	Channel uint8
}

func MakeMIDIChannelPrefix(tick uint64, delta lib.Delta, channel uint8) MIDIChannelPrefix {
	if channel > 15 {
		panic(fmt.Sprintf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel))
	}

	return MIDIChannelPrefix{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x20, 0x01, channel}),
			tag:    lib.TagMIDIChannelPrefix,
			Status: 0xff,
			Type:   lib.TypeMIDIChannelPrefix,
		},
		Channel: channel,
	}
}

func UnmarshalMIDIChannelPrefix(tick uint64, delta lib.Delta, bytes []byte) (*MIDIChannelPrefix, error) {
	if len(bytes) != 1 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix length (%d): expected '1'", len(bytes))
	}

	if channel := bytes[0]; channel > 15 {
		return nil, fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	} else {
		event := MakeMIDIChannelPrefix(tick, delta, channel)

		return &event, nil
	}
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
	m.tag = lib.TagMIDIChannelPrefix
	m.Status = 0xff
	m.Type = lib.TypeMIDIChannelPrefix

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)MIDIChannelPrefix\s+([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid MIDIChannelPrefix event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if channel < 0 || channel > 15 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	} else {
		m.delta = delta
		m.Channel = uint8(channel)
	}

	return nil
}

func (e MIDIChannelPrefix) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag     string    `json:"tag"`
		Delta   lib.Delta `json:"delta"`
		Status  byte      `json:"status"`
		Type    byte      `json:"type"`
		Channel uint8     `json:"channel"`
	}{
		Tag:     fmt.Sprintf("%v", e.tag),
		Delta:   e.delta,
		Status:  byte(e.Status),
		Type:    byte(e.Type),
		Channel: e.Channel,
	}

	return json.Marshal(t)
}

func (e *MIDIChannelPrefix) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag     string    `json:"tag"`
		Delta   lib.Delta `json:"delta"`
		Channel uint8     `json:"channel"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagMIDIChannelPrefix) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagMIDIChannelPrefix
		e.Type = lib.TypeMIDIChannelPrefix
		e.Channel = t.Channel
	}

	return nil
}
