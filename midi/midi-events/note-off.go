package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type NoteOff struct {
	MidiEvent
	Note     byte
	Velocity byte
}

func NewNoteOff(event MidiEvent, r *bufio.Reader) (*NoteOff, error) {
	if event.Status&0xF0 != 0x80 {
		return nil, fmt.Errorf("Invalid NoteOff status (%02x): expected '80'", event.Status&0xF0)
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

	return &NoteOff{
		MidiEvent: event,
		Note:      note,
		Velocity:  velocity,
	}, nil
}

func (e *NoteOff) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d note:%d velocity:%d\n", e.MidiEvent, "NoteOff", e.Channel, e.Note, e.Velocity)
}
