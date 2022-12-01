package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
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
			Status:  lib.Status(0xd0 | channel),
			Channel: channel,
		},
		Pressure: pressure,
	}
}

func UnmarshalChannelPressure(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte) (*ChannelPressure, error) {
	if status&0xf0 != 0xd0 {
		return nil, fmt.Errorf("Invalid ChannelPressure status (%v): expected 'Dx'", status)
	}

	if len(data) < 1 {
		return nil, fmt.Errorf("Invalid ChannelPressure data (%v): expected pressure", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var pressure = data[0]

	if channel > 15 {
		return nil, fmt.Errorf("invalid channel (%v)", channel)
	}

	if pressure > 127 {
		return nil, fmt.Errorf("invalid pressure (%v)", pressure)
	}

	event := MakeChannelPressure(tick, delta, channel, pressure)

	return &event, nil
}

func (p ChannelPressure) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xd0 | p.Channel),
		p.Pressure,
	}

	return
}

func (p *ChannelPressure) UnmarshalText(bytes []byte) error {
	p.tick = 0
	p.delta = 0
	p.bytes = []byte{}
	p.tag = lib.TagChannelPressure
	p.Status = 0xd0

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)ChannelPressure\s+channel:([0-9]+)\s+pressure:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
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
		p.delta = delta
		p.Status = or(p.Status, channel)
		p.Channel = channel
		p.Pressure = uint8(pressure)
	}

	return nil
}
