package metaevent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/transcriptaze/midiasm/midi/types"
)

type EndOfTrack struct {
	event
}

func MakeEndOfTrack(tick uint64, delta uint32) EndOfTrack {
	return EndOfTrack{
		event: event{
			tick:   tick,
			delta:  delta,
			bytes:  []byte{0x00, 0xff, 0x2f, 0x00},
			tag:    types.TagEndOfTrack,
			Status: 0xff,
			Type:   types.TypeEndOfTrack,
		},
	}
}

func UnmarshalEndOfTrack(tick uint64, delta uint32, bytes []byte) (*EndOfTrack, error) {
	event := MakeEndOfTrack(tick, delta)

	return &event, nil
}

func (e EndOfTrack) MarshalBinary() (encoded []byte, err error) {
	return []byte{
		byte(e.Status),
		byte(e.Type),
		byte(0),
	}, nil
}

func (e *EndOfTrack) UnmarshalText(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.tag = types.TagEndOfTrack
	e.Type = 0x2f

	re := regexp.MustCompile(`(?i)delta:([0-9]+)(?:.*?)EndOfTrack`)
	text := string(bytes)

	if match := re.FindStringSubmatch(text); match == nil || len(match) < 2 {
		return fmt.Errorf("invalid EndOfTrack event (%v)", text)
	} else if delta, err := strconv.ParseUint(match[1], 10, 32); err != nil {
		return err
	} else {
		e.delta = uint32(delta)
	}

	return nil
}

func (e EndOfTrack) MarshalJSON() (encoded []byte, err error) {
	t := struct {
		Tag    string `json:"tag"`
		Delta  uint32 `json:"delta"`
		Status byte   `json:"status"`
		Type   byte   `json:"type"`
	}{
		Tag:    fmt.Sprintf("%v", e.tag),
		Delta:  e.delta,
		Status: byte(e.Status),
		Type:   byte(e.Type),
	}

	return json.Marshal(t)
}

func (e *EndOfTrack) UnmarshalJSON(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.tag = types.TagEndOfTrack
	e.Type = 0x2f

	t := struct {
		Tag   string `json:"tag"`
		Delta uint32 `json:"delta"`
	}{}

	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	} else if t.Tag != "EndOfTrack" {
		return fmt.Errorf("invalid EndOfTrack event (%v)", string(bytes))
	} else {
		e.delta = t.Delta
	}

	return nil
}
