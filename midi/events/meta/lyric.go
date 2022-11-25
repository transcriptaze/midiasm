package metaevent

import (
	"encoding/json"
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

func (e *Lyric) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Lyric\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Lyric event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagLyric
		e.Status = 0xff
		e.Type = lib.TypeLyric
		e.Lyric = match[2]
	}

	return nil
}

func (e Lyric) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Type   byte      `json:"type"`
		Lyric  string    `json:"lyric"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
		Lyric:  e.Lyric,
	}

	return json.Marshal(t)
}

func (e *Lyric) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
		Lyric string    `json:"lyric"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagLyric) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagLyric
		e.Type = lib.TypeLyric
		e.Lyric = t.Lyric
	}

	return nil
}
