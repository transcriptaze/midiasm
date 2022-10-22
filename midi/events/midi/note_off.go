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

func NewNoteOff(tick uint64, delta uint32, channel uint8, note Note, velocity uint8, bytes ...byte) *NoteOff {
	return &NoteOff{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: bytes,

			Tag:     "NoteOff",
			Status:  types.Status(0x80 | channel),
			Channel: types.Channel(channel),
		},
		Note:     note,
		Velocity: velocity,
	}
}

func UnmarshalNoteOff(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (*NoteOff, error) {
	if status&0xf0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%v): expected '8x'", status)
	}

	channel := types.Channel(status & 0x0f)

	note, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	velocity, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &NoteOff{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: r.Bytes(),

			Tag:     "NoteOff",
			Status:  status,
			Channel: channel,
		},
		Note: Note{
			Value: note,
			Name:  GetNoteOff(ctx, channel, note),
			Alias: FormatNote(ctx, note),
		},
		Velocity: velocity,
	}, nil
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
	n.Tag = "NoteOff"

	re := regexp.MustCompile(`(?i)NoteOff\s+channel:([0-9]+)\s+note:([A-G][♯♭]?[0-9]),\s*velocity:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 4 {
		return fmt.Errorf("invalid NoteOff event (%v)", text)
	} else if channel, err := strconv.ParseUint(match[1], 10, 8); err != nil {
		return err
	} else if note, err := ParseNote(nil, match[2]); err != nil {
		return err
	} else if velocity, err := strconv.ParseUint(match[3], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid NoteOff channel (%v)", channel)
	} else if velocity > 127 {
		return fmt.Errorf("invalid NoteOff velocity (%v)", channel)
	} else {
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
