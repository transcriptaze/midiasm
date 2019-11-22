package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type NoteOn struct {
	MidiEvent
	Note     byte
	Velocity byte
}

func NewNoteOn(event MidiEvent, r *bufio.Reader) (*NoteOn, error) {
	if event.Status&0xF0 != 0x90 {
		return nil, fmt.Errorf("Invalid NoteOn status (%02x): expected '90'", event.Status&0x80)
	}

	note, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	velocity, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	event.bytes = append(event.bytes, note, velocity)

	return &NoteOn{
		MidiEvent: event,
		Note:      note,
		Velocity:  velocity,
	}, nil
}

func (e *NoteOn) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d note:%d velocity:%d\n", e.MidiEvent, "NoteOn", e.Channel, e.Note, e.Velocity)
}
