package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type DeviceName struct {
	MetaEvent
	Name string
}

func NewDeviceName(event *MetaEvent, r io.ByteReader) (*DeviceName, error) {
	if event.Type != 0x09 {
		return nil, fmt.Errorf("Invalid DeviceName event type (%02x): expected '09'", event.Type)
	}

	name, err := read(r)
	if err != nil {
		return nil, err
	}

	return &DeviceName{
		MetaEvent: *event,
		Name:      string(name),
	}, nil
}

func (e *DeviceName) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "DeviceName", e.Name)
}
