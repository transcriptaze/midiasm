package midievent

import (
	"bufio"
	"fmt"
	"io"
)

type ProgramChange struct {
	MidiEvent
	Program byte
}

func NewProgramChange(event MidiEvent, r *bufio.Reader) (*ProgramChange, error) {
	if event.Status&0xF0 != 0xc0 {
		return nil, fmt.Errorf("Invalid ProgramChange status (%02x): expected 'C0'", event.Status&0x80)
	}

	program, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	event.bytes = append(event.bytes, program)

	return &ProgramChange{
		MidiEvent: event,
		Program:   program,
	}, nil
}

func (e *ProgramChange) Render(w io.Writer) {
	fmt.Fprintf(w, "   ")
	for i := 5; i > len(e.bytes); i-- {
		fmt.Fprintf(w, "   ")
	}
	for _, b := range e.bytes {
		fmt.Fprintf(w, "%02x ", b)
	}
	fmt.Fprintf(w, "                                     ")

	fmt.Fprintf(w, "%02x/%-16s %s channel:%d program:%d\n", e.Status, "ProgramChange", e.MidiEvent.Event, e.Channel, e.Program)
}
