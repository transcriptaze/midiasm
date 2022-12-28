package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

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

func (e *Marker) unmarshal(tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	marker := string(data)

	*e = MakeMarker(tick, delta, marker, bytes...)

	return nil
}

func (m Marker) MarshalBinary() (encoded []byte, err error) {
	return append([]byte{
		byte(m.Status),
		byte(m.Type),
		byte(len(m.Marker)),
	},
		[]byte(m.Marker)...), nil
}

func (e *Marker) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagMarker, remaining[0])
	} else if !equals(remaining[1], lib.TypeMarker) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagMarker, remaining[1])
	} else if marker, err := vlf(remaining[2:]); err != nil {
		return err
	} else {
		*e = MakeMarker(0, delta, string(marker), bytes...)
	}

	return nil
}

func (e *Marker) UnmarshalText(bytes []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)Marker\s+(.*)`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 3 {
		return fmt.Errorf("invalid Marker event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeMarker(0, delta, match[2], []byte{}...)
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
		*e = MakeMarker(0, t.Delta, t.Marker, []byte{}...)
	}

	return nil
}
