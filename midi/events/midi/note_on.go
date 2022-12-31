package midievent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type NoteOn struct {
	event
	Note     Note
	Velocity byte
}

func MakeNoteOn(tick uint64, delta uint32, channel lib.Channel, note Note, velocity uint8, bytes ...byte) NoteOn {
	if channel > 15 {
		panic(fmt.Sprintf("invalid channel (%v)", channel))
	}

	if velocity > 127 {
		panic(fmt.Sprintf("invalid velocity (%v)", velocity))
	}

	return NoteOn{
		event: event{
			tick:  tick,
			delta: lib.Delta(delta),
			bytes: bytes,

			tag:     lib.TagNoteOn,
			Status:  lib.Status(0x90 | channel),
			Channel: channel,
		},
		Note:     note,
		Velocity: velocity,
	}
}

func (e *NoteOn) unmarshal(tick uint64, delta uint32, status lib.Status, data []byte, bytes ...byte) error {
	if status&0xf0 != 0x90 {
		return fmt.Errorf("Invalid NoteOn status (%v): expected '9x'", status)
	}

	if len(data) < 2 {
		return fmt.Errorf("Invalid NoteOff data (%v): expected note and velocity", data)
	}

	var channel = lib.Channel(status & 0x0f)
	var note = Note{
		Value: data[0],
		Name:  FormatNote(nil, data[0]),
		Alias: FormatNote(nil, data[0]),
	}
	var velocity uint8

	if v := data[1]; v > 127 {
		return fmt.Errorf("Invalid NoteOn velocity (%v)", v)
	} else {
		velocity = v
	}

	*e = MakeNoteOn(tick, delta, channel, note, velocity, bytes...)

	return nil
}

func (n NoteOn) MarshalBinary() (encoded []byte, err error) {
	encoded = []byte{
		byte(0x90 | n.Channel),
		n.Note.Value,
		n.Velocity,
	}

	return
}

func (e *NoteOn) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) != 3 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if !lib.TypeNoteOn.Equals(remaining[0]) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagNoteOn, remaining[0])
	} else if channel := remaining[0] & 0x0f; channel > 15 {
		return fmt.Errorf("invalid channel (%v)", channel)
	} else if data := remaining[1:]; len(data) < 2 {
		return fmt.Errorf("Invalid NoteOn data")
	} else if velocity := data[1]; velocity > 127 {
		return fmt.Errorf("Invalid NoteOn velocity (%v)", velocity)
	} else {
		note := Note{
			Value: data[0],
			Name:  FormatNote(nil, data[0]),
			Alias: FormatNote(nil, data[0]),
		}

		*e = MakeNoteOn(0, delta, lib.Channel(channel), note, velocity, bytes...)
	}

	return nil
}

func (e *NoteOn) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)NoteOn\s+channel:([0-9]+)\s+note:([A-G][♯♭]?[0-9]),\s*velocity:([0-9]+)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 5 {
		return fmt.Errorf("invalid NoteOn event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else if channel, err := lib.ParseChannel(match[2]); err != nil {
		return err
	} else if note, err := ParseNote(nil, match[3]); err != nil {
		return err
	} else if velocity, err := strconv.ParseUint(match[4], 10, 8); err != nil {
		return err
	} else if velocity > 127 {
		return fmt.Errorf("invalid NoteOn velocity (%v)", velocity)
	} else {
		*e = MakeNoteOn(0, uint32(delta), lib.Channel(channel), note, uint8(velocity), []byte{}...)
	}

	return nil
}

func (e NoteOn) MarshalJSON() (encoded []byte, err error) {
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

func (e *NoteOn) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag      string      `json:"tag"`
		Delta    lib.Delta   `json:"delta"`
		Channel  lib.Channel `json:"channel"`
		Note     Note        `json:"note"`
		Velocity uint8       `json:"velocity"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagNoteOn) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		note := Note{
			Value: t.Note.Value,
			Name:  FormatNote(nil, t.Note.Value),
			Alias: FormatNote(nil, t.Note.Value),
		}

		*e = MakeNoteOn(0, uint32(t.Delta), t.Channel, note, t.Velocity, []byte{}...)
	}

	return nil
}

func (e NoteOn) Transpose(ctx *context.Context, steps int) NoteOn {
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

	return NoteOn{
		event: event{
			tick:    e.tick,
			delta:   e.delta,
			tag:     lib.TagNoteOn,
			Status:  lib.Status(0x90 | e.Channel),
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

func (e NoteOn) Format(ctx *context.Context) NoteOn {
	return NoteOn{
		event: event{
			tick:    e.tick,
			delta:   e.delta,
			bytes:   e.bytes,
			tag:     e.tag,
			Status:  e.Status,
			Channel: e.Channel,
		},
		Note: Note{
			Value: e.Note.Value,
			Name:  ctx.FormatNote(e.Note.Value),
			Alias: ctx.FormatNote(e.Note.Value),
		},
		Velocity: e.Velocity,
	}
}

func ParseNote(ctx *context.Context, s string) (Note, error) {
	re := regexp.MustCompile(`([A-G][♯♭]?)([-]?[0-9])`)

	if match := re.FindStringSubmatch(s); match == nil || len(match) != 3 {
		return Note{}, fmt.Errorf("invalid note (%v)", s)
	} else if note, err := lib.LookupNote(match[1]); err != nil {
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
