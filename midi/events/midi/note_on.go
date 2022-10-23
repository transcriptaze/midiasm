package midievent

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/io"
	"github.com/transcriptaze/midiasm/midi/types"
)

type NoteOn struct {
	event
	Note     Note
	Velocity byte
}

func NewNoteOn(ctx *context.Context, tick uint64, delta uint32, r IO.Reader, status types.Status) (*NoteOn, error) {
	if status&0xf0 != 0x90 {
		return nil, fmt.Errorf("Invalid NoteOn status (%v): expected '9x'", status)
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

	formatted := FormatNote(ctx, note)
	if ctx != nil {
		ctx.PutNoteOn(channel, note)
	}

	return &NoteOn{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: r.Bytes(),

			tag:     types.TagNoteOn,
			Status:  status,
			Channel: channel,
		},
		Note: Note{
			Value: note,
			Name:  formatted,
			Alias: formatted,
		},
		Velocity: velocity,
	}, nil
}

func (n *NoteOn) Transpose(ctx *context.Context, steps int) {
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
		Name:  FormatNote(ctx, note),
		Alias: FormatNote(ctx, note),
	}
}

func (n NoteOn) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0x90 | n.Channel),
		n.Note.Value,
		n.Velocity,
	}

	return
}

func (n *NoteOn) UnmarshalText(bytes []byte) error {
	n.tick = 0
	n.delta = 0
	n.bytes = []byte{}
	n.tag = types.TagNoteOn

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)NoteOn\s+channel:([0-9]+)\s+note:([A-G][♯♭]?[0-9]),\s*velocity:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid NoteOn event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else if channel, err := strconv.ParseUint(match[2], 10, 8); err != nil {
		return err
	} else if note, err := ParseNote(nil, match[3]); err != nil {
		return err
	} else if velocity, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if channel > 15 {
		return fmt.Errorf("invalid NoteOn channel (%v)", channel)
	} else if velocity > 127 {
		return fmt.Errorf("invalid NoteOn velocity (%v)", channel)
	} else {
		n.delta = uint32(delta)
		n.bytes = []byte{0x00, byte(0x90 | uint8(channel&0x0f)), note.Value, byte(velocity)}
		n.Status = types.Status(0x90 | uint8(channel&0x0f))
		n.Channel = types.Channel(channel)
		n.Note = note
		n.Velocity = uint8(velocity)
	}

	return nil
}

func FormatNote(ctx *context.Context, n byte) string {
	if ctx != nil {
		return ctx.FormatNote(n)
	}

	var scale = context.Sharps
	var note = scale[n%12]
	var octave int

	if context.MiddleC == types.C4 {
		octave = int(n/12) - 2
	} else {
		octave = int(n/12) - 1
	}

	return fmt.Sprintf("%s%d", note, octave)
}

func ParseNote(ctx *context.Context, s string) (Note, error) {
	re := regexp.MustCompile(`([A-G][♯♭]?)([-]?[0-9])`)

	if match := re.FindStringSubmatch(s); match == nil || len(match) != 3 {
		return Note{}, fmt.Errorf("invalid note (%v)", s)
	} else if note, err := types.LookupNote(match[1]); err != nil {
		return Note{}, err
	} else if octave, err := strconv.ParseInt(match[2], 10, 8); err != nil {
		return Note{}, err
	} else {
		return Note{
			Value: uint8(12*int8(octave+1) + int8(note.Ord)),
			Name:  fmt.Sprintf("%s%d", note.Name, octave),
			Alias: fmt.Sprintf("%s%d", note.Name, octave),
		}, nil
	}
}
