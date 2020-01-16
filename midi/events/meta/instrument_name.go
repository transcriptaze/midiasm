package metaevent

import (
	"fmt"
	"io"
)

type InstrumentName struct {
	MetaEvent
	Name string
}

func NewInstrumentName(event *MetaEvent, r io.ByteReader) (*InstrumentName, error) {
	if event.Type != 0x04 {
		return nil, fmt.Errorf("Invalid InstrumentName event type (%02x): expected '04'", event.Type)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &InstrumentName{
		MetaEvent: *event,
		Name:      string(name),
	}, nil
}

func (e *InstrumentName) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "InstrumentName", e.Name)
}
