package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	lib "github.com/transcriptaze/midiasm/midi/types"
)

type NoteOff struct {
	event
	Note     Note
	Velocity byte
}

func MakeNoteOff(tick uint64, delta uint32, channel lib.Channel, note Note, velocity uint8, bytes ...byte) NoteOff {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	if velocity > 127 {
		panic(fmt.Sprintf("invalid velocity (%v)", velocity))
	}

	return NoteOff{
		event: event{
			tick:    tick,
			delta:   lib.Delta(delta),
			bytes:   bytes,
			tag:     lib.TagNoteOff,
			Status:  lib.Status(0x80 | channel),
			Channel: channel,
		},
		Note:     note,
		Velocity: velocity,
	}
}

func UnmarshalNoteOff(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte) (*NoteOff, error) {
	if status&0xf0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%v): expected '8x'", status)
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("Invalid NoteOff data (%v): expected note and velocity", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var note = Note{
		Value: data[0],
		Name:  GetNoteOff(ctx, channel, data[0]),
		Alias: FormatNote(ctx, data[0]),
	}
	var velocity uint8

	if v := data[1]; v > 127 {
		return nil, fmt.Errorf("Invalid NoteOff velocity (%v)", v)
	} else {
		velocity = v
	}

	event := MakeNoteOff(tick, delta, channel, note, velocity)

	return &event, nil
}

func (e NoteOff) Transpose(ctx *context.Context, steps int) NoteOff {
	v := int(e.Note.Value) + steps
	note := e.Note.Value

	switch {
	case v < 0:
		note = 0

	case v > 127:
		note = 127

	default:
		note = byte(v)
	}

	return NoteOff{
		event: event{
			tick:    e.tick,
			delta:   e.delta,
			bytes:   []byte{},
			tag:     lib.TagNoteOff,
			Status:  lib.Status(0x80 | e.Channel),
			Channel: e.Channel,
		},
		Note: Note{
			Value: note,
			Name:  ctx.GetNoteOff(e.Channel, note),
			Alias: ctx.FormatNote(note),
		},
		Velocity: e.Velocity,
	}
}

func (n NoteOff) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0x80 | n.Channel),
		n.Note.Value,
		n.Velocity,
	}

	return
}

func (n *NoteOff) UnmarshalText(bytes []byte) error {
	n.tick = 0
	n.delta = 0
	n.bytes = []byte{}
	n.tag = lib.TagNoteOff
	n.Status = 0x80

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)NoteOff\s+channel:([0-9]+)\s+note:([A-G][♯♭]?[0-9]),\s*velocity:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid NoteOff event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if note, err := ParseNote(nil, match[3]); err != nil {
		return err
	} else if velocity, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if velocity > 127 {
		return fmt.Errorf("invalid NoteOff velocity (%v)", channel)
	} else {
		n.delta = delta
		n.Status = or(n.Status, channel)
		n.Channel = channel
		n.Note = note
		n.Velocity = uint8(velocity)
	}

	return nil
}

func GetNoteOff(ctx *context.Context, ch lib.Channel, n byte) string {
	if ctx != nil {
		return ctx.GetNoteOff(ch, n)
	}

	return FormatNote(ctx, n)
}
