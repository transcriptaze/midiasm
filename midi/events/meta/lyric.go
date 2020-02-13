package metaevent

import (
	"github.com/twystd/midiasm/midi/events"
	"github.com/twystd/midiasm/midi/types"
	"io"
)

type Lyric struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Lyric  string
}

func NewLyric(r io.ByteReader) (*Lyric, error) {
	lyric, err := events.VLF(r)
	if err != nil {
		return nil, err
	}

	return &Lyric{
		Tag:    "Lyric",
		Status: 0xff,
		Type:   0x05,
		Lyric:  string(lyric),
	}, nil
}
