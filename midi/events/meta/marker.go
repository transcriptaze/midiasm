package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type Marker struct {
	event
	Marker string
}

func MakeMarker(tick uint64, delta lib.Delta, marker string, bytes ...byte) Marker {
	return Marker{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagMarker,
			Status: 0xff,
			Type:   lib.TypeMarker,
		},
		Marker: marker,
	}
}

func UnmarshalMarker(ctx *context.Context, tick uint64, delta lib.Delta, data []byte, bytes ...byte) (*Marker, error) {
	marker := string(data)
	event := MakeMarker(tick, delta, marker, bytes...)

	return &event, nil
}

func (m Marker) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(m.Status),
		byte(m.Type),
		byte(len(m.Marker)),
	},
		[]byte(m.Marker)...), nil
}

func (e *Marker) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Marker\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Marker event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		e.tick = 0
		e.delta = delta
		e.bytes = []byte{}
		e.tag = lib.TagMarker
		e.Status = 0xff
		e.Type = lib.TypeMarker
		e.Marker = match[2]
	}

	return nil
}

func (e Marker) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Type   byte      `json:"type"`
		Marker string    `json:"marker"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
		Marker: e.Marker,
	}

	return json.Marshal(t)
}

func (e *Marker) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Marker string    `json:"marker"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if !equal(t.Tag, lib.TagMarker) {
		return fmt.Errorf("invalid %v event (%v)", e.tag, string(bytes))
	} else {
		e.tick = 0
		e.delta = t.Delta
		e.bytes = []byte{}
		e.Status = 0xff
		e.tag = lib.TagMarker
		e.Type = lib.TypeMarker
		e.Marker = t.Marker
	}

	return nil
}
