package metaevent

import (
	"fmt"
	"regexp"
	"strconv"
)

type EndOfTrack struct {
	event
}

func NewEndOfTrack(tick uint64, delta uint32, bytes []byte) (*EndOfTrack, error) {
	if len(bytes) != 0 {
		return nil, fmt.Errorf("Invalid EndOfTrack length (%d): expected '0'", len(bytes))
	}

	return &EndOfTrack{
		event: event{
			tick:  tick,
			delta: delta,
			bytes: []byte{0x00, 0xff, 0x2f, 0x00},

			Tag:    "EndOfTrack",
			Status: 0xff,
			Type:   0x2f,
		},
	}, nil
}

func (e EndOfTrack) MarshalBinary() (encoded []byte, err error) {
	encoded = make([]byte, 3)

	encoded[0] = byte(e.Status)
	encoded[1] = byte(e.Type)
	encoded[2] = byte(0)

	return
}

func (e *EndOfTrack) UnmarshalText(bytes []byte) error {
	e.tick = 0
	e.delta = 0
	e.bytes = []byte{}
	e.Status = 0xff
	e.Tag = "EndOfTrack"
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
