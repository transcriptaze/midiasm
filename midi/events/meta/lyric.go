package metaevent

import (
	"github.com/transcriptaze/midiasm/midi/types"
)

type Lyric struct {
	Tag    string
	Status types.Status
	Type   types.MetaEventType
	Lyric  string
}

func NewLyric(bytes []byte) (*Lyric, error) {
	return &Lyric{
		Tag:    "Lyric",
		Status: 0xff,
		Type:   0x05,
		Lyric:  string(bytes),
	}, nil
}
