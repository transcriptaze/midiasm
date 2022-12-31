package midievent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type ChannelPressure struct {
	event
	Pressure byte
}

func MakeChannelPressure(tick uint64, delta uint32, channel lib.Channel, pressure uint8, bytes ...byte) ChannelPressure {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	if pressure > 127 {
		panic(fmt.Sprintf("invalid pressure (%v)", pressure))
	}

	return ChannelPressure{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagChannelPressure,
			Status:  or(0xD0, channel),
			Channel: channel,
		},
		Pressure: pressure,
	}
}

func (e *ChannelPressure) unmarshal(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0xd0 {
		return fmt.Errorf("Invalid ChannelPressure status (%v): expected 'Dx'", status)
	}

	if len(data) < 1 {
		return fmt.Errorf("Invalid ChannelPressure data (%v): expected pressure", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var pressure = data[0]

	if channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	}

	if pressure > 127 {
		return fmt.Errorf("invalid pressure (%v)", pressure)
	}

	*e = MakeChannelPressure(tick, delta, channel, pressure, bytes...)

	return nil
}

func (e ChannelPressure) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xD0 | e.Channel),
		e.Pressure,
	}

	return
}

func (e *ChannelPressure) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) != 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeChannelPressure.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagChannelPressure, remaining[0])
	} else if channel := remaining[0] & 0x0f; channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	} else if data := remaining[1:]; len(data) < 1 {
		return fmt.Errorf("Invalid ChannelPressure data")
	} else if pressure := data[0]; pressure > 127 {
		return fmt.Errorf("InvalidChannelPressure pressure (%v)", pressure)
	} else {
		*e = MakeChannelPressure(0, delta, lib.Channel(channel), pressure, bytes...)
	}

	return nil
}

func (e *ChannelPressure) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)ChannelPressure\s+channel:([0-9]+)\s+pressure:([0-9]+)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid ChannelPressure event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if pressure, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if pressure > 127 {
		return fmt.Errorf("invalid ChannelPressure pressure (%v)", pressure)
	} else {
		*e = MakeChannelPressure(0, uint32(delta), lib.Channel(channel), uint8(pressure), []byte{}...)
	}

	return nil
}

func (e ChannelPressure) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag      string      `json:"tag"`
		Delta    lib.Delta   `json:"delta"`
		Status   byte        `json:"status"`
		Channel  lib.Channel `json:"channel"`
		Pressure uint8       `json:"pressure"`
	}{
		Tag:      fmt.Sprintf("%v", e.tag),
		Delta:    e.delta,
		Status:   byte(e.Status),
		Channel:  e.Channel,
		Pressure: e.Pressure,
	}

	return json.Marshal(t)
}

func (e *ChannelPressure) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag      string      `json:"tag"`
		Delta    lib.Delta   `json:"delta"`
		Channel  lib.Channel `json:"channel"`
		Pressure uint8       `json:"pressure"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagChannelPressure) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		*e = MakeChannelPressure(0, uint32(t.Delta), t.Channel, t.Pressure, []byte{}...)
	}

	return nil
}
