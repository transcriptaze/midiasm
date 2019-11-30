package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/event"
	"io"
)

type Tempo struct {
	MetaEvent
	Tempo uint32
}

func NewTempo(event *MetaEvent, r io.ByteReader) (*Tempo, error) {
	if event.Type != 0x51 {
		return nil, fmt.Errorf("Invalid Tempo event type (%02x): expected '51'", event.Type)
	}

	data, err := read(r)
	if err != nil {
		return nil, err
	} else if len(data) != 3 {
		return nil, fmt.Errorf("Invalid Tempo length (%d): expected '3'", len(data))
	}

	tempo := uint32(0)
	for _, b := range data {
		tempo <<= 8
		tempo += uint32(b)
	}

	return &Tempo{
		MetaEvent: *event,
		Tempo:     tempo,
	}, nil
}

func (e *Tempo) Render(ctx *event.Context, w io.Writer) {
	fmt.Fprintf(w, "%s %-16s tempo:%v", e.MetaEvent, "Tempo", e.Tempo)
}
