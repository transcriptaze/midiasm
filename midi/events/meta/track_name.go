package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type TrackName struct {
	MetaEvent
	Name string
}

func NewTrackName(event *MetaEvent, r io.ByteReader) (*TrackName, error) {
	if event.Type != 0x03 {
		return nil, fmt.Errorf("Invalid TrackName event type (%02x): expected '03'", event.Type)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	event.Tag = "TrackName"
	return &TrackName{
		MetaEvent: *event,
		Name:      string(name),
	}, nil
}

func (e *TrackName) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "TrackName", e.Name)
}
