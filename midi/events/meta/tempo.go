package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
)

type Tempo struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Tempo  uint32
}

func NewTempo(bytes []byte) (*Tempo, error) {
	if len(bytes) != 3 {
		return nil, fmt.Errorf("Invalid Tempo length (%d): expected '3'", len(bytes))
	}

	tempo := uint32(0)
	for _, b := range bytes {
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
