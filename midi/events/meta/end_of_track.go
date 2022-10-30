package metaevent

import (
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
