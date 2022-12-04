package midievent

import (
	"encoding/json"
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

func (e *PitchBend) unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0xe0 {
		return fmt.Errorf("Invalid PitchBend status (%v): expected 'Ex'", status)
	}

	if len(data) < 2 {
		return fmt.Errorf("Invalid PitchBend data (%v): expected bend", data)
	}

	channel := lib.Channel(status & 0x0f)
	bend := uint16(data[0])
	bend <<= 7
	bend |= uint16(data[1]) & 0x7f

	if channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	}

	*e = MakePitchBend(tick, delta, channel, bend, bytes...)

	return nil
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
func (e PitchBend) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag     string      `json:"tag"`
		Delta   lib.Delta   `json:"delta"`
		Status  byte        `json:"status"`
		Channel lib.Channel `json:"channel"`
		Bend    uint16      `json:"bend"`
	}{
		Tag:     fmt.Sprintf("%v", e.tag),
		Delta:   e.delta,
		Status:  byte(e.Status),
		Channel: e.Channel,
		Bend:    e.Bend,
	}

	return json.Marshal(t)
}

func (e *PitchBend) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag     string      `json:"tag"`
		Delta   lib.Delta   `json:"delta"`
		Channel lib.Channel `json:"channel"`
		Bend    uint16      `json:"bend"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagPitchBend) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.tag = lib.TagPitchBend
		e.Status = or(0xE0, t.Channel)
		e.Channel = t.Channel
		e.Bend = t.Bend
	}

	return nil
}
