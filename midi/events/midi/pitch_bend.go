package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type PitchBend struct {
	event
	Bend uint16
}

func MakePitchBend(tick uint64, delta uint32, channel lib.Channel, bend uint16, bytes ...byte) PitchBend {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	return PitchBend{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagPitchBend,
			Status:  lib.Status(0xe0 | channel),
			Channel: channel,
		},
		Bend: bend,
	}
}

func UnmarshalPitchBend(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) (*PitchBend, error) {
	if status&0xf0 != 0xe0 {
		return nil, fmt.Errorf("Invalid PitchBend status (%v): expected 'Ex'", status)
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("Invalid PitchBend data (%v): expected bend", data)
	}

	channel := lib.Channel(status & 0x0f)
	bend := uint16(data[0])
	bend <<= 7
	bend |= uint16(data[1]) & 0x7f

	if channel > 15 {
		return nil, fmt.Errorf("invalid channel (%v)", channel)
	}

	event := MakePitchBend(tick, delta, channel, bend, bytes...)

	return &event, nil
}

func (b PitchBend) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0xe0 | b.Channel),
		byte(b.Bend >> 8 & 0x00ff),
		byte(b.Bend >> 0 & 0x00ff),
	}

	return
}

func (b *PitchBend) UnmarshalText(bytes []byte) error {
	b.tick = 0
	b.delta = 0
	b.bytes = []byte{}
	b.tag = lib.TagPitchBend
	b.Status = 0xe0

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)PitchBend\s+channel:([0-9]+)\s+bend:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid PitchBend event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if bend, err := strconv.ParseUint(match[3], 10, 16); err != nil {
		return err
	} else {
		b.delta = delta
		b.Status = or(b.Status, channel)
		b.Channel = channel
		b.Bend = uint16(bend)
	}

	return nil
}
