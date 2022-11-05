package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type NoteOff struct {
	event
	Note     Note
	Velocity byte
}

func MakeNoteOff(tick uint64, delta uint32, channel types.Channel, note Note, velocity uint8, bytes ...byte) NoteOff {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	if velocity > 127 {
		panic(fmt.Sprintf("invalid velocity (%v)", velocity))
	}

	return NoteOff{
		event: event{
			tick:  tick,
			delta: types.Delta(delta),
			bytes: bytes,

			tag:     types.TagNoteOff,
			Status:  types.Status(0x80 | channel),
			Channel: channel,
		},
		Note:     note,
		Velocity: velocity,
	}
}

func UnmarshalNoteOff(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (*NoteOff, error) {
	if status&0xf0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%v): expected '8x'", status)
	}

	var channel = types.Channel(status & 0x0f)
	var note Note
	var velocity uint8

	if n, err := r.ReadByte(); err != nil {
		return nil, err
	} else {
		note.Value = n
		note.Name = GetNoteOff(ctx, channel, n)
		note.Alias = FormatNote(ctx, n)
	}

	if v, err := r.ReadByte(); err != nil {
		return nil, err
	} else if v > 127 {
		return nil, fmt.Errorf("invalid NoteOn velocity (%v)", v)
	} else {
		velocity = v
	}

	event := MakeNoteOff(tick, delta, channel, note, velocity, r.Bytes()...)

	return &event, nil
}

func (n *NoteOff) Transpose(ctx *context.Context, steps int) {
	v := int(n.Note.Value) + steps
	note := n.Note.Value

	switch {
	case v < 0:
		note = 0

	case v > 127:
		note = 127

	default:
		note = byte(v)
	}

	n.Note = Note{
		Value: note,
		Name:  ctx.GetNoteOff(n.Channel, note),
		Alias: FormatNote(ctx, note),
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
	n.tag = types.TagNoteOff

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)NoteOff\s+channel:([0-9]+)\s+note:([A-G][♯♭]?[0-9]),\s*velocity:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid NoteOff event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if note, err := ParseNote(nil, match[3]); err != nil {
		return err
	} else if velocity, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid NoteOff channel (%v)", channel)
	} else if velocity > 127 {
		return fmt.Errorf("invalid NoteOff velocity (%v)", channel)
	} else {
		n.delta = types.Delta(delta)
		n.bytes = []byte{0x00, byte(0x80 | uint8(channel&0x0f)), note.Value, byte(velocity)}
		n.Status = types.Status(0x80 | uint8(channel&0x0f))
		n.Channel = types.Channel(channel)
		n.Note = note
		n.Velocity = uint8(velocity)
	}

	return nil
}

func GetNoteOff(ctx *context.Context, ch types.Channel, n byte) string {
	if ctx != nil {
		return ctx.GetNoteOff(ch, n)
	}

	return FormatNote(ctx, n)
}
