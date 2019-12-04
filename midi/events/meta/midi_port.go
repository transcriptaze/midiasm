package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/context"
	"io"
)

type MIDIPort struct {
	MetaEvent
	Port uint8
}

func NewMIDIPort(event *MetaEvent, r io.ByteReader) (*MIDIPort, error) {
	if event.Type != 0x21 {
		return nil, fmt.Errorf("Invalid MIDIPort event type (%02x): expected '21'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 1 {
		return nil, fmt.Errorf("Invalid MIDIPort length (%d): expected '1'", len(data))
	}

	port := data[0]
	if port < 0 || port > 127 {
		return nil, fmt.Errorf("Invalid MIDIPort port (%d): expected a value in the interval [0..127]", port)
	}

	return &MIDIPort{
		MetaEvent: *event,
		Port:      port,
	}, nil
}

func (e *MIDIPort) Render(ctx *context.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s %d", e.MetaEvent, "MIDIPort", e.Port)
}
