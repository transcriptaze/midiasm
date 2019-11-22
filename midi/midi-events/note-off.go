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
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d note:%d velocity:%d\n", e.Status, "NoteOff", e.MidiEvent.Event, e.Channel, e.Note, e.Velocity)
}
