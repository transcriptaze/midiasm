package midievent

import (
	"fmt"
	"io"
)

type NoteOff struct {
	MidiEvent
	Note     byte
	Velocity byte
}

func NewNoteOff(event *MidiEvent, r io.ByteReader) (*NoteOff, error) {
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

	return &NoteOff{
		MidiEvent: *event,
		Note:      note,
		Velocity:  velocity,
	}, nil
}

func (e *NoteOff) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d note:%d velocity:%d", e.MidiEvent, "NoteOff", e.Channel, e.Note, e.Velocity)
}
