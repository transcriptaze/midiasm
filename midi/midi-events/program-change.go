package midievent

import (
	"fmt"
	"io"
)

type ProgramChange struct {
	MidiEvent
	Program byte
}

func NewProgramChange(event *MidiEvent, r io.ByteReader) (*ProgramChange, error) {
	if event.Status&0xF0 != 0xc0 {
		return nil, fmt.Errorf("Invalid ProgramChange status (%02x): expected 'C0'", event.Status&0x80)
	}

	program, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	return &ProgramChange{
		MidiEvent: *event,
		Program:   program,
	}, nil
}

func (e *ProgramChange) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s channel:%d program:%d", e.MidiEvent, "ProgramChange", e.Channel, e.Program)
}
