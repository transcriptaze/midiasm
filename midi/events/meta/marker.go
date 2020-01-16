package metaevent

import (
	"fmt"
	"io"
)

type Marker struct {
	MetaEvent
	Marker string
}

func NewMarker(event *MetaEvent, r io.ByteReader) (*Marker, error) {
	if event.Type != 0x06 {
		return nil, fmt.Errorf("Invalid Marker event type (%02x): expected '06'", event.Type)
	}

	marker, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Marker{
		MetaEvent: *event,
		Marker:    string(marker),
	}, nil
}

func (e *Marker) Render(w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "Marker", e.Marker)
}
