package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type PitchBend struct {
	event
	Bend uint16
}

func MakePitchBend(tick uint64, delta uint32, channel types.Channel, bend uint16, bytes ...byte) PitchBend {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	return PitchBend{
		event: event{
			tick:    tick,
			delta:   types.Delta(delta),
			bytes:   bytes,
			tag:     types.TagPitchBend,
			Status:  types.Status(0xe0 | channel),
			Channel: channel,
		},
		Bend: bend,
	}
}

func UnmarshalPitchBend(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (*PitchBend, error) {
	if status&0xf0 != 0xe0 {
		return nil, fmt.Errorf("Invalid PitchBend status (%v): expected 'Ex'", status)
	}

	channel := types.Channel(status & 0x0f)
	bend := uint16(0)

	for i := 0; i < 2; i++ {
		if b, err := r.ReadByte(); err != nil {
			return nil, err
		} else {
			bend <<= 7
			bend |= uint16(b) & 0x7f
		}
	}

	event := MakePitchBend(tick, delta, channel, bend, r.Bytes()...)

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
	b.tag = types.TagPitchBend

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)PitchBend\s+channel:([0-9]+)\s+bend:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid PitchBend event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid PitchBend channel (%v)", channel)
	} else if bend, err := strconv.ParseUint(match[3], 10, 16); err != nil {
		return err
	} else {
		b.delta = types.Delta(delta)
		b.bytes = []byte{}
		b.Status = types.Status(0xe0 | uint8(channel&0x0f))
		b.Channel = types.Channel(channel)
		b.Bend = uint16(bend)

		if bytes, err := b.delta.MarshalBinary(); err != nil {
			return err
		} else {
			b.bytes = append(b.bytes, bytes...)
			if bytes, err = b.MarshalBinary(); err != nil {
				return err
			} else {
				b.bytes = append(b.bytes, bytes...)
			}
		}
	}

	return nil
}
