package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Tempo struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Tempo  uint32
}

func NewTempo(r io.ByteReader) (*Tempo, error) {
	data, err := events.VLF(r)
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
		Tag:    "Tempo",
		Status: 0xff,
		Type:   0x51,
		Tempo:  tempo,
	}, nil
}
