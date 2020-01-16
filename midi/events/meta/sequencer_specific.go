package metaevent

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type SequencerSpecificEvent struct {
	MetaEvent
	Data []byte
}

func NewSequencerSpecificEvent(event *MetaEvent, r io.ByteReader) (*SequencerSpecificEvent, error) {
	if event.Type != 0x7f {
		return nil, fmt.Errorf("Invalid SequencerSpecificEvent event type (%02x): expected '7F'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	}

	return &SequencerSpecificEvent{
		MetaEvent: *event,
		Data:      data,
	}, nil
}

func (e *SequencerSpecificEvent) Render(w io.Writer) {
	data := new(bytes.Buffer)

	for _, b := range e.Data {
		fmt.Fprintf(data, "%02X ", b)
	}

	fmt.Fprintf(w, "%s %-16s %s", e.MetaEvent, "SequencerSpecificEvent", strings.TrimSpace(data.String()))
}
