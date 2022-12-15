package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type MIDIChannelPrefix struct {
	event
	Channel uint8
}

func MakeMIDIChannelPrefix(tick uint64, delta lib.Delta, channel uint8, bytes ...byte) MIDIChannelPrefix {
	if channel > 15 {
		panic(fmt.Sprintf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel))
	}

	return MIDIChannelPrefix{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagMIDIChannelPrefix,
			Status: 0xff,
			Type:   lib.TypeMIDIChannelPrefix,
		},
		Channel: channel,
	}
}

func (e *MIDIChannelPrefix) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	if len(data) != 1 {
		return fmt.Errorf("Invalid MIDIChannelPrefix length (%d): expected '1'", len(data))
	}

	if channel := data[0]; channel > 15 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	} else {
		*e = MakeMIDIChannelPrefix(tick, delta, channel, bytes...)

		return nil
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

func (e *MIDIChannelPrefix) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagMIDIChannelPrefix, remaining[0])
	} else if !equals(remaining[1], lib.TypeMIDIChannelPrefix) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagMIDIChannelPrefix, remaining[1])
	} else if data, err := vlf(remaining[2:]); err != nil {
		return err
	} else if len(data) < 1 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel data")
	} else if channel := data[0]; channel > 15 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	} else {
		*e = MakeMIDIChannelPrefix(0, delta, channel, bytes...)
	}

	return nil
}

func (e *MIDIChannelPrefix) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)MIDIChannelPrefix\s+([0-9]+)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid MIDIChannelPrefix event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", channel)
	} else {
		*e = MakeMIDIChannelPrefix(0, delta, uint8(channel), []byte{}...)
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
	} else if t.Channel > 15 {
		return fmt.Errorf("Invalid MIDIChannelPrefix channel (%d): expected a value in the interval [0..15]", t.Channel)
	} else {
		*e = MakeMIDIChannelPrefix(0, t.Delta, t.Channel, []byte{}...)
	}

	return nil
}
