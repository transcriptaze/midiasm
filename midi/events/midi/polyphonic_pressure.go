package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type PolyphonicPressure struct {
	event
	Pressure byte
}

func MakePolyphonicPressure(tick uint64, delta uint32, channel types.Channel, pressure uint8) PolyphonicPressure {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	if pressure > 127 {
		panic(fmt.Sprintf("invalid pressure (%v)", pressure))
	}

	return PolyphonicPressure{
		event: event{
			tick:    tick,
			delta:   types.Delta(delta),
			bytes:   []byte{0x00, byte(0xa0 | channel), pressure},
			tag:     types.TagPolyphonicPressure,
			Status:  types.Status(0xa0 | channel),
			Channel: channel,
		},
		Pressure: pressure,
	}
}

func UnmarshalPolyphonicPressure(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (*PolyphonicPressure, error) {
	if status&0xf0 != 0xa0 {
		return nil, fmt.Errorf("Invalid PolyphonicPressure status (%v): expected 'Ax'", status)
	}

	var channel = types.Channel(status & 0x0f)
	var pressure uint8

	if v, err := r.ReadByte(); err != nil {
		return nil, err
	} else if v > 127 {
		return nil, fmt.Errorf("invalid PolyphonicPressure pressure (%v)", v)
	} else {
		pressure = v
	}

	event := MakePolyphonicPressure(tick, delta, channel, pressure)

	return &event, nil
}

func (p PolyphonicPressure) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xa0 | p.Channel),
		p.Pressure,
	}

	return
}

func (p *PolyphonicPressure) UnmarshalText(bytes []byte) error {
	p.tick = 0
	p.delta = 0
	p.bytes = []byte{}
	p.tag = types.TagPolyphonicPressure

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)PolyphonicPressure\s+channel:([0-9]+)\s+pressure:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid PolyphonicPressure event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid PolyphonicPressure channel (%v)", channel)
	} else if pressure, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if pressure > 127 {
		return fmt.Errorf("invalid PolyphonicPressure pressure (%v)", pressure)
	} else {
		p.delta = types.Delta(delta)
		p.bytes = []byte{0x00, byte(0xa0 | uint8(channel&0x0f)), byte(pressure)}
		p.Status = types.Status(0xa0 | uint8(channel&0x0f))
		p.Channel = types.Channel(channel)
		p.Pressure = uint8(pressure)
	}

	return nil
}
