package metaevent

import (
	"fmt"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Lyric struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Lyric  string
}

func NewLyric(r io.ByteReader, status types.Status, eventType types.MetaEventType) (*Lyric, error) {
	if eventType != 0x05 {
		return nil, fmt.Errorf("Invalid Lyric event type (%02x): expected '05'", eventType)
	}

	lyric, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Lyric{
		Tag:    "Lyric",
		Status: status,
		Type:   eventType,
		Lyric:  string(lyric),
	}, nil
}
