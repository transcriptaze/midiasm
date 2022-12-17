package midievent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
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

func (e *NoteOff) unmarshal(ctx *context.Context, tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0x80 {
		return fmt.Errorf("Invalid NoteOff status (%v): expected '8x'", status)
	}

	if len(data) < 2 {
		return fmt.Errorf("Invalid NoteOff data (%v): expected note and velocity", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var note = Note{
		Value: data[0],
		Name:  GetNoteOff(ctx, channel, data[0]),
		Alias: FormatNote(ctx, data[0]),
	}
	var velocity uint8

	if v := data[1]; v > 127 {
		return fmt.Errorf("Invalid NoteOff velocity (%v)", v)
	} else {
		velocity = v
	}

	*e = MakeNoteOff(tick, delta, channel, note, velocity, bytes...)

	return nil
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

func (e *NoteOff) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) != 3 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeNoteOff.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagNoteOff, remaining[0])
	} else if channel := remaining[0] & 0x0f; channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	} else if data := remaining[1:]; len(data) < 2 {
		return fmt.Errorf("Invalid NoteOff data")
	} else if velocity := data[1]; velocity > 127 {
		return fmt.Errorf("Invalid NoteOff velocity (%v)", velocity)
	} else {
		note := Note{
			Value: data[0],
			Name:  GetNoteOff(nil, lib.Channel(channel), data[0]),
			Alias: FormatNote(nil, data[0]),
		}

		*e = MakeNoteOff(0, delta, lib.Channel(channel), note, velocity, bytes...)
	}

	return nil
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

func (e NoteOff) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag      string      `json:"tag"`
		Delta    lib.Delta   `json:"delta"`
		Status   byte        `json:"status"`
		Channel  lib.Channel `json:"channel"`
		Note     Note        `json:"note"`
		Velocity uint8       `json:"velocity"`
	}{
		Tag:      fmt.Sprintf("%v", e.tag),
		Delta:    e.delta,
		Status:   byte(e.Status),
		Channel:  e.Channel,
		Note:     e.Note,
		Velocity: e.Velocity,
	}

	return json.Marshal(t)
}

func (e *NoteOff) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag      string      `json:"tag"`
		Delta    lib.Delta   `json:"delta"`
		Channel  lib.Channel `json:"channel"`
		Note     Note        `json:"note"`
		Velocity uint8       `json:"velocity"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagNoteOff) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.tag = lib.TagNoteOff
		e.Status = or(0x80, t.Channel)
		e.Channel = t.Channel
		e.Note = Note{
			Value: t.Note.Value,
			Name:  FormatNote(nil, t.Note.Value),
			Alias: FormatNote(nil, t.Note.Value),
		}
		e.Velocity = t.Velocity
	}

	return nil
}

func GetNoteOff(ctx *context.Context, ch lib.Channel, n byte) string {
	if ctx != nil {
		return ctx.GetNoteOff(ch, n)
	}

	return FormatNote(ctx, n)
}
