package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/transcriptaze/midiasm/midi/context"
	"github.com/transcriptaze/midiasm/midi/lib"
)

type EndOfTrack struct {
	event
}

func MakeEndOfTrack(tick uint64, delta lib.Delta, bytes ...byte) EndOfTrack {
	return EndOfTrack{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  bytes,
			tag:    lib.TagEndOfTrack,
			Status: 0xff,
			Type:   lib.TypeEndOfTrack,
		},
	}
}

func (e *EndOfTrack) unmarshal(ctx *context.Context, tick uint64, delta lib.Delta, status byte, data []byte, bytes ...byte) error {
	*e = MakeEndOfTrack(tick, delta, bytes...)

	return nil
}

func (e EndOfTrack) MarshalBinary() (encoded []byte, err error) {
	return []byte{
		byte(e.Status),
		byte(e.Type),
		byte(0),
	}, nil
}

func (e *EndOfTrack) UnmarshalBinary(bytes []byte) error {
	if delta, remaining, err := delta(bytes); err != nil {
		return err
	} else if len(remaining) < 2 {
		return fmt.Errorf("Invalid event (%v)", remaining)
	} else if remaining[0] != 0xff {
		return fmt.Errorf("Invalid %v status (%02X)", lib.TagEndOfTrack, remaining[0])
	} else if !equals(remaining[1], lib.TypeEndOfTrack) {
		return fmt.Errorf("Invalid %v event type (%02X)", lib.TagEndOfTrack, remaining[1])
	} else {
		*e = MakeEndOfTrack(0, delta, bytes...)
	}

	return nil
}

func (e *EndOfTrack) UnmarshalText(text []byte) error {
	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)EndOfTrack`)

	if match := re.FindStringSubmatch(string(text)); match == nil || len(match) < 2 {
		return fmt.Errorf("invalid EndOfTrack event (%v)", text)
	} else if delta, err := lib.ParseDelta(match[1]); err != nil {
		return err
	} else {
		*e = MakeEndOfTrack(0, delta, []byte{}...)
	}

	return nil
}

func (e EndOfTrack) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string    `json:"tag"`
		Delta  lib.Delta `json:"delta"`
		Status byte      `json:"status"`
		Type   byte      `json:"type"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
	}

	return json.Marshal(t)
}

func (e *EndOfTrack) UnmarshalJSON(bytes []byte) error {
	t := struct {
		Tag   string    `json:"tag"`
		Delta lib.Delta `json:"delta"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if t.Tag != "EndOfTrack" {
		return fmt.Errorf("invalid EndOfTrack event (%v)", string(bytes))
	} else {
		*e = MakeEndOfTrack(0, t.Delta, []byte{}...)
	}

	return nil
}
