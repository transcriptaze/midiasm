package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type ProgramName struct {
	MetaEvent
	Name string
}

func NewProgramName(event *MetaEvent, r io.ByteReader) (*ProgramName, error) {
	if event.Type != 0x08 {
		return nil, fmt.Errorf("Invalid ProgramName event type (%02x): expected '08'", event.Type)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &ProgramName{
		MetaEvent: *event,
		Name:      string(name),
	}, nil
}

func (e *ProgramName) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "ProgramName", e.Name)
}
