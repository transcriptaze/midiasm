package metaevent

import (
	"fmt"
	"regexp"

	lib "github.com/transcriptaze/midiasm/midi/types"
)

type Lyric struct {
	event
	Lyric string
}

func MakeLyric(tick uint64, delta lib.Delta, lyric string) Lyric {
	return Lyric{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  append([]byte{0x00, 0xff, 0x05, byte(len(lyric))}, []byte(lyric)...),
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   lib.TypeLyric,
		},
		Lyric: lyric,
	}
}

func UnmarshalLyric(tick uint64, delta lib.Delta, bytes []byte) (*Lyric, error) {
	lyric := string(bytes)
	event := MakeLyric(tick, delta, lyric)

	return &event, nil
}

func (l Lyric) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(l.Status),
		byte(l.Type),
		byte(len(l.Lyric)),
	},
		[]byte(l.Lyric)...), nil
}

func (l *Lyric) UnmarshalText(bytes []byte) error {
	l.tick = 0
	l.delta = 0
	l.bytes = []byte{}
	l.tag = lib.TagLyric
	l.Status = 0xff
	l.Type = lib.TypeLyric

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Lyric\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Lyric event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		l.delta = delta
		l.Lyric = string(match[2])
	}

	return nil
}
