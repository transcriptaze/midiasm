package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/lib"
)

type Lyric struct {
	event
	Lyric string
}

func MakeLyric(tick uint64, delta lib.Delta, lyric string, bytes ...byte) Lyric {
	return Lyric{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagLyric,
			Status: 0xff,
			Type:   lib.TypeLyric,
		},
		Lyric: lyric,
	}
}

func (e *Lyric) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	lyric := string(data)

	*e = MakeLyric(tick, delta, lyric, bytes...)

	return nil
}

func (l Lyric) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(l.Status),
		byte(l.Type),
		byte(len(l.Lyric)),
	},
		[]byte(l.Lyric)...), nil
}

func (e *Lyric) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagLyric, remaining[0])
	} else if !equals(remaining[1], lib.TypeLyric) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagLyric, remaining[1])
	} else if lyric, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeLyric(0, delta, string(lyric), bytes...)
	}

	return nil
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
