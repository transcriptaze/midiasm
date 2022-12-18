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

func (e *PolyphonicPressure) unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0xa0 {
		return fmt.Errorf("Invalid PolyphonicPressure status (%v): expected 'Ax'", status)
	}

	if len(data) < 1 {
		return fmt.Errorf("Invalid PolyphonicPressure data (%v): expected pressure", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var pressure = data[0]

	*e = MakePolyphonicPressure(tick, delta, channel, pressure, bytes...)

	return nil
}

func (e PolyphonicPressure) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xA0 | e.Channel),
		e.Pressure,
	}

	return
}

func (e *PolyphonicPressure) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) != 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypePolyphonicPressure.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagPolyphonicPressure, remaining[0])
	} else if channel := remaining[0] & 0x0f; channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	} else if data := remaining[1:]; len(data) < 1 {
		return fmt.Errorf("Invalid PolyphonicPressure data")
	} else if pressure := data[0]; pressure > 127 {
		return fmt.Errorf("Invalid PolyphonicPressure pressure (%v)", pressure)
	} else {
		*e = MakePolyphonicPressure(0, delta, lib.Channel(channel), pressure, bytes...)
	}

	return nil
}

func (e *PolyphonicPressure) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)PolyphonicPressure\s+channel:([0-9]+)\s+pressure:([0-9]+)`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 4 {
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
		*e = MakePolyphonicPressure(0, uint32(delta), lib.Channel(channel), uint8(pressure), []byte{}...)
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
		*e = MakePolyphonicPressure(0, uint32(t.Delta), t.Channel, t.Pressure, []byte{}...)
	}

	return nil
}
