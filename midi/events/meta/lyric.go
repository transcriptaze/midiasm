package metaevent

import (
	"fmt"
	"io"
)

type Lyric struct {
	MetaEvent
	Lyric string
}

func NewLyric(event *MetaEvent, r io.ByteReader) (*Lyric, error) {
	if event.Type != 0x05 {
		return nil, fmt.Errorf("Invalid Lyric event type (%02x): expected '05'", event.Type)
	}

	lyric, err := read(r)
	if err != nil {
		return nil, err
	}

	return &Lyric{
		MetaEvent: *event,
		Lyric:     string(lyric),
	}, nil
}
