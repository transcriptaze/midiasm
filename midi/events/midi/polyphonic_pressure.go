package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	lib "github.com/transcriptaze/midiasm/midi/types"
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
			Status:  lib.Status(0xa0 | channel),
			Channel: channel,
		},
		Pressure: pressure,
	}
}

func UnmarshalPolyphonicPressure(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (*PolyphonicPressure, error) {
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
	p.tag = lib.TagPolyphonicPressure
	p.Status = 0xa0

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
		p.delta = delta
		p.Status = or(p.Status, channel)
		p.Channel = channel
		p.Pressure = uint8(pressure)
	}

	return nil
}
