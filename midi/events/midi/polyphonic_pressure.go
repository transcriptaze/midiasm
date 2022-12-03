package midievent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type PolyphonicPressure struct {
	event
	Pressure byte
}

func MakePolyphonicPressure(tick uint64, delta uint32, channel lib.Channel, pressure uint8, bytes ...byte) PolyphonicPressure {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	if pressure > 127 {
		panic(fmt.Sprintf("invalid pressure (%v)", pressure))
	}

	return PolyphonicPressure{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagPolyphonicPressure,
			Status:  or(0xA0, channel),
			Channel: channel,
		},
		Pressure: pressure,
	}
}

func UnmarshalPolyphonicPressure(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (*PolyphonicPressure, error) {
	if status&0xf0 != 0xa0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%v): expected 'Ax'", status)
	}

	if len(data) < 1 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure data (%v): expected pressure", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var pressure = data[0]

	event := MakePolyphonicPressure(tick, delta, channel, pressure, bytes...)

	return &event, nil
}

func (e PolyphonicPressure) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xA0 | e.Channel),
		e.Pressure,
	}

	return
}

func (e *PolyphonicPressure) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)PolyphonicPressure\s+channel:([0-9]+)\s+pressure:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid PolyphonicPressure event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if pressure, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if pressure > 127 {
		return fmt.Errorf("invalid PolyphonicPressure pressure (%v)", pressure)
	} else {
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagPolyphonicPressure
		e.Status = or(0xA0, channel)
		e.Channel = channel
		e.Pressure = uint8(pressure)
	}

	return nil
}

func (e PolyphonicPressure) MarshalJSON() (encoded []byte, err error) {
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

func (e *PolyphonicPressure) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag      string      `json:"tag"`
		Delta    lib.Delta   `json:"delta"`
		Channel  lib.Channel `json:"channel"`
		Pressure uint8       `json:"pressure"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagPolyphonicPressure) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.tag = lib.TagPolyphonicPressure
		e.Status = lib.Status(0xA0 | t.Channel)
		e.Channel = t.Channel
		e.Pressure = t.Pressure
	}

	return nil
}
